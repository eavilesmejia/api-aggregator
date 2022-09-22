package restful

import (
    "context"
    internalTx "github.com/INFURA/infra/internal/transaction"
    "github.com/INFURA/infra/repository/infura"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/gofiber/fiber/v2"
    "net/http"
)

/**
 * Router setup + route handlers
 */

type transaction struct {
}

func (t *transaction) SetupRoutes(router fiber.Router) {
    router.Get("/transactions/:txHash", getTransaction)
    router.Post("/transactions", searchTransaction)
}

func newTransactionRouter() *transaction {
    return &transaction{}
}

/**
 * API Models
 */

type jsonErr struct {
    JsonCode int    `json:"jsonCode,omitempty"`
    HttpCode uint   `json:"httpCode,omitempty"`
    Message  string `json:"error,omitempty"`
}

type getTransactionResponse struct {
    Block uint64 `json:"block,omitempty"`
    Value uint64 `json:"value,omitempty"`
    Price uint64 `json:"price,omitempty"`
    Gas   uint64 `json:"gas,omitempty"`
    Hash  string `json:"hash,omitempty"`
    From  string `json:"from,omitempty"`
    To    string `json:"to,omitempty"`
}

type searchTransactionRequest struct {
    BlockString string `json:"blockString,omitempty"`
    BlockNumber int64  `json:"blockNumber,omitempty"`
    TxIndex     int64  `json:"txIndex,omitempty"` // use value 1 for first index position
    Hash        string `json:"hash,omitempty"`
}

/**
 * API Functions
 */

// getTransaction retrieves a transaction by its hash
func getTransaction(c *fiber.Ctx) error {
    txHash := c.Params("txHash")
    if txHash == "" {
        msgErr := jsonErr{
            HttpCode: http.StatusBadRequest,
            Message:  "transaction hash is empty",
        }
        return c.Status(http.StatusBadRequest).JSON(msgErr)
    }
    resp, msgErr := makeAndParseTransaction(c.Context(), "", 0, -1, txHash)
    if msgErr != nil {
        return c.Status(http.StatusBadRequest).JSON(msgErr)
    }
    
    return c.Status(http.StatusOK).JSON(resp)
}

// searchTransaction retrieves transaction by "searchTransactionRequest" fields
func searchTransaction(c *fiber.Ctx) error {
    req := new(searchTransactionRequest)
    if err := c.BodyParser(req); err != nil {
        return err
    }
    req.TxIndex -= 1
    resp, msgErr := makeAndParseTransaction(c.Context(), req.BlockString, req.BlockNumber, req.TxIndex, req.Hash)
    if msgErr != nil {
        return c.Status(http.StatusBadRequest).JSON(msgErr)
    }
    
    return c.Status(http.StatusOK).JSON(resp)
}

func makeAndParseTransaction(ctx context.Context, blockString string, blockNumber int64, txIndex int64, hash string) (*getTransactionResponse, *jsonErr) {
    repository := infura.NewInfura()
    tx, err := internalTx.GetTransaction(ctx, repository, blockString, blockNumber, txIndex, hash)
    if err != nil {
        msgErr := &jsonErr{
            HttpCode: http.StatusInternalServerError,
            Message:  err.Error(),
        }
        return nil, msgErr
    }
    resp := &getTransactionResponse{
        Hash: tx.Hash,
        From: tx.From,
        To:   tx.To,
    }
    resp.Block, _ = hexutil.DecodeUint64(tx.BlockNumber)
    resp.Value, _ = hexutil.DecodeUint64(tx.Value)
    resp.Gas, _ = hexutil.DecodeUint64(tx.Gas)
    resp.Price, _ = hexutil.DecodeUint64(tx.GasPrice)
    return resp, nil
}
