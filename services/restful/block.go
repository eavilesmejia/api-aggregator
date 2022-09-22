package restful

import (
    "github.com/INFURA/infra/internal/block"
    "github.com/INFURA/infra/repository/infura"
    "github.com/gofiber/fiber/v2"
    "net/http"
)

/**
 * Router setup + route handlers
 */

type blockRouter struct {
}

func (b *blockRouter) SetupRoutes(router fiber.Router) {
    router.Post("/blocks", getBlock)
}

func newBlockRouter() *blockRouter {
    return &blockRouter{}
}

/**
 * API Models
 */

/**
 * API Functions
 */

type getBlockRequest struct {
    Number int64  `json:"number,omitempty"`
    Hash   string `json:"hash,omitempty"`
}

type getBlockResponse struct {
    TransactionAmount int    `json:"transactionAmount,omitempty"`
    Number            uint64 `json:"number,omitempty"`
    Difficulty        uint64 `json:"difficulty,omitempty"`
    Hash              string `json:"hash"`
    UncleHash         string `json:"uncleHash,omitempty"`
    ParentHash        string `json:"parentHash,omitempty"`
}

func getBlock(c *fiber.Ctx) error {
    req := new(getBlockRequest)
    if err := c.BodyParser(req); err != nil {
        return err
    }
    repository := infura.NewInfura()
    bloq, err := block.GetBlock(c.Context(), repository, req.Number, req.Hash, false)
    if err != nil {
        msgErr := &jsonErr{
            HttpCode: http.StatusInternalServerError,
            Message:  err.Error(),
        }
        return c.Status(http.StatusInternalServerError).JSON(msgErr)
    }
    resp := &getBlockResponse{
        Number:            bloq.Block.NumberU64(),
        Hash:              bloq.Block.Hash().String(),
        UncleHash:         bloq.Block.UncleHash().String(),
        ParentHash:        bloq.Block.ParentHash().String(),
        Difficulty:        bloq.Block.Difficulty().Uint64(),
        TransactionAmount: len(bloq.Block.Transactions()),
    }
    return c.Status(http.StatusOK).JSON(resp)
}
