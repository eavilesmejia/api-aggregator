package infura

import (
    "context"
    "github.com/INFURA/infra/internal"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
)

func TestMakeRequest(t *testing.T) {
    t.Run("eth_getTransactionByBlockNumberAndIndex", func(t *testing.T) {
        var result internal.Transaction
        msgErr, err := makeRequest(context.Background(), &result, "eth_getTransactionByBlockNumberAndIndex", "0x5BAD55", "0x0")
        if err != nil {
            t.Fatal(err)
        }
        if msgErr != nil {
            t.Fatalf("error code: %d, msg: %s", msgErr.Code, msgErr.Message)
        }
        if result.Hash != "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f" {
            t.Errorf("expected value: %s, got: %s", "0x8784d99762bccd03b2086eabccee0d77f14d05463281e121a62abfebcf0d2d5f", result.Hash)
        }
    })
    
    t.Run("eth_getBlockByNumber", func(t *testing.T) {
        type customBlock struct {
            types.Header
            types.Transactions
        }
        var result customBlock
        msgErr, err := makeRequest(context.Background(), &result, "eth_getBlockByNumber", "0x5BAD55", false)
        if err != nil {
            t.Fatal(err)
        }
        if msgErr != nil {
            t.Fatalf("error code: %d, msg: %s", msgErr.Code, msgErr.Message)
        }
        if result.Number.Int64() != int64(6008149) {
            t.Errorf("expected value: %d, got: %d", 6008149, result.Number.Int64())
        }
    })
    
    t.Run("eth_getTransactionByHash", func(t *testing.T) {
        var result internal.Transaction
        msgErr, err := makeRequest(context.Background(), &result, "eth_getTransactionByHash", common.HexToHash("0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0"))
        if err != nil {
            t.Fatal(err)
        }
        if msgErr != nil {
            t.Fatalf("error code: %d, msg: %s", msgErr.Code, msgErr.Message)
        }
        if result.Hash != "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0" {
            t.Errorf("expected value: %s, got: %s", "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0", result.Hash)
        }
    })
    t.Run("Non-Blocking", func(t *testing.T) {
        var result string
        // we want to hit local interface to avoid networking delays
        fake1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Println("req1: hello")
            time.Sleep(time.Second * 5)
            w.WriteHeader(http.StatusOK)
            w.Header().Set("Content-Type", "application/json")
            log.Println("req1: world")
            _, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
        }))
        fake2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Println("req2: hello")
            time.Sleep(time.Second * 2)
            w.WriteHeader(http.StatusOK)
            w.Header().Set("Content-Type", "application/json")
            log.Println("req2: world")
            _, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
        }))
        _ = os.Setenv("INFURA_API_URL", fake1.URL)
        
        jErr, err := makeRequest(context.Background(), &result, "eth_chainId")
        if err != nil {
            t.Fatal(err)
        }
        if jErr != nil {
            t.Fatalf(jErr.Message)
        }
        
        _ = os.Setenv("INFURA_API_URL", fake2.URL)
        jErr, err = makeRequest(context.Background(), &result, "eth_chainId")
        if err != nil {
            t.Fatal(err)
        }
        if jErr != nil {
            t.Fatalf(jErr.Message)
        }
    })
}

func BenchmarkMakeRequest(b *testing.B) {
    
    // we want to hit local interface to avoid networking delays
    fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        
        w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
    }))
    _ = os.Setenv("INFURA_API_URL", fake.URL)
    
    for i := 0; i < b.N; i++ {
        b.Run("eth_chainId", func(b *testing.B) {
            var result string
            _, err := makeRequest(context.Background(), &result, "eth_chainId")
            if err != nil {
                b.Fatal(err)
            }
        })
    }
}

func BenchmarkClient(b *testing.B) {
    // in case we want to hit local interface to avoid networking delays
    fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        
        w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
    }))
    _ = os.Setenv("INFURA_API_URL", fake.URL)
    
    client, err := getEthClient(context.Background())
    if err != nil {
        b.Fatal(err)
    }
    for i := 0; i < b.N; i++ {
        b.Run("eth_chainId", func(b *testing.B) {
            _, err := client.ChainID(context.Background())
            if err != nil {
                b.Fatal(err)
            }
        })
    }
}
