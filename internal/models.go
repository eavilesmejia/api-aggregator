package internal

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	Block *types.Block
	Err   Error
}

type BqParams struct {
	Hash        string
	BlockNumber int64
	Details     bool
}

type Transaction struct {
	Gas         string
	GasPrice    string
	Value       string
	Type        string
	BlockHash   string
	BlockNumber string
	From        string
	Hash        string
	Input       string
	Nonce       string
	R           string
	S           string
	To          string
	IsPending   bool
	Err         Error
}

type TxParams struct {
	Hash        string
	BlockNumber int64
	TxIndex     int64
	BlockString string
}

type Error struct {
	Code    int
	Message string
}

func (e *Error) String() string {
	return fmt.Sprintf("error msg: %s, error code: %d", e.Message, e.Code)
}
