package infura

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/ethereum/go-ethereum/ethclient"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"
)

var transporter *http.Transport

const (
    id              = 1
    version         = "2.0"
    responseTimeout = time.Second * 10
)

type response struct {
    data []byte
    err  error
}

type jsonError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type jsonrpcMessage struct {
    Version string          `json:"jsonrpc,omitempty"`
    ID      json.RawMessage `json:"id,omitempty"`
    Method  string          `json:"method,omitempty"`
    Params  json.RawMessage `json:"params,omitempty"`
    Error   *jsonError      `json:"error,omitempty"`
    Result  json.RawMessage `json:"result,omitempty"`
}

func init() {
    // declare a transporter globally to reuse connection
    transporter = &http.Transport{}
}

// doRequest sends an HTTP request to infura URL.
//
// Checks for a valid http code response
func doRequest(ctx context.Context, msg *jsonrpcMessage, responseChan chan response) {
    u, err := url.Parse(os.Getenv("INFURA_API_URL"))
    if err != nil {
        responseChan <- response{data: nil, err: err}
        return
    }
    jsonPayload, err := json.Marshal(msg)
    if err != nil {
        responseChan <- response{
            data: nil,
            err:  err,
        }
        return
    }
    
    // declare client with a reusable connection
    client := &http.Client{Transport: transporter}
    
    jsonPayloadStr := string(jsonPayload)
    body := strings.NewReader(jsonPayloadStr)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
    
    if err != nil {
        responseChan <- response{
            data: nil,
            err:  err,
        }
        return
    }
    
    resp, err := client.Do(req)
    if err != nil {
        responseChan <- response{
            data: nil,
            err:  err,
        }
        return
    }
    
    defer resp.Body.Close()
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        responseChan <- response{data: nil, err: err}
        return
    }
    
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        responseChan <- response{data: nil, err: fmt.Errorf("invalid http status response %d", resp.StatusCode)}
        return
    }
    
    responseChan <- response{data: bodyBytes, err: nil}
}

// prepareRequest marshalls to the proper JSON-RPC format that infura waits.
func prepareRequest(method string, params ...interface{}) (*jsonrpcMessage, error) {
    msg := &jsonrpcMessage{
        Version: version,
        ID:      strconv.AppendUint(nil, uint64(id), 10),
        Method:  method,
    }
    if params != nil { // prevent sending "params":null
        var err error
        if msg.Params, err = json.Marshal(params); err != nil {
            return nil, err
        }
    }
    return msg, nil
}

// parserResponse unmarshalls JSON-RPC response into a given pointer.
// If "result" field is nil or not a pointer it will not unmarshall the response.
func parserResponse(response response, result interface{}) (*jsonError, error) {
    if response.err != nil {
        return nil, response.err
    }
    
    var msg jsonrpcMessage
    
    if err := json.Unmarshal(response.data, &msg); err != nil {
        return nil, err
    }
    if msg.Error != nil {
        return msg.Error, nil
    }
    
    if err := json.Unmarshal(msg.Result, &result); err != nil {
        return nil, err
    }
    return nil, nil
}

// makeRequest makes a JSON-RPC call to the given method and unmarshall it.
// If "result" field is nil or not a pointer it will not unmarshall the response.
func makeRequest(ctx context.Context, result interface{}, method string, params ...interface{}) (*jsonError, error) {
    // we use an unbuffered channel since we are declaring it for each request
    responseChan := make(chan response)
    payload, err := prepareRequest(method, params...)
    if err != nil {
        return nil, err
    }
    go doRequest(ctx, payload, responseChan)
    defer close(responseChan)
    
    var jErr *jsonError
    timer := time.NewTimer(responseTimeout)
    done := make(chan bool)
    
    go func() {
        select {
        // there was not any http calling error, so let's parse the response
        case response := <-responseChan:
            jErr, err = parserResponse(response, result)
            close(done)
        // if the context was cancelled for any reason
        case <-ctx.Done():
            jErr, err = nil, ctx.Err()
            close(done)
        // define a timeout to avoid a waiting request too long
        case <-timer.C:
            jErr, err = nil, fmt.Errorf("response timeout")
            close(done)
        }
    }()
    <-done
    return jErr, err
}

func getEthClient(ctx context.Context) (*ethclient.Client, error) {
    rawUrl, err := url.Parse(os.Getenv("INFURA_API_URL"))
    if err != nil {
        return nil, err
    }
    
    return ethclient.DialContext(context.Background(), rawUrl.String())
}
