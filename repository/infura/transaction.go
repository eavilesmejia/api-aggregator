package infura

import (
	"context"
	model "github.com/INFURA/infra/internal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type transaction struct {
}

func (t *transaction) getTransactionByHash(ctx context.Context, hash string) (*model.Transaction, error) {
	var result model.Transaction
	jErr, err := makeRequest(ctx, &result, "eth_getTransactionByHash", common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	if jErr != nil {
		result.Err.Message = jErr.Message
		result.Err.Code = jErr.Code
	}
	return &result, nil
}

func (t *transaction) getTransactionByNumberAndIndex(ctx context.Context, blockString string, number int64, index int64) (*model.Transaction, error) {
	var result model.Transaction
	var jErr *jsonError
	var err error
	indexHex := hexutil.EncodeBig(big.NewInt(index))
	if blockString == "latest" ||
		blockString == "earliest" ||
		blockString == "pending" {
		jErr, err = makeRequest(ctx, &result, "eth_getTransactionByBlockNumberAndIndex", blockString, indexHex)
	} else {
		numberHex := hexutil.EncodeBig(big.NewInt(number))
		jErr, err = makeRequest(ctx, &result, "eth_getTransactionByBlockNumberAndIndex", numberHex, indexHex)
	}

	if err != nil {
		return nil, err
	}

	if jErr != nil {
		result.Err.Message = jErr.Message
		result.Err.Code = jErr.Code
	}

	if blockString == "pending" {
		result.IsPending = true
	}

	return &result, nil
}

func (t *transaction) getTransactionByBlockHashAndIndex(ctx context.Context, hash string, index int64) (*model.Transaction, error) {
	var result model.Transaction
	indexHex := hexutil.EncodeBig(big.NewInt(index))
	jErr, err := makeRequest(ctx, &result, "eth_getTransactionByBlockHashAndIndex", common.HexToHash(hash), indexHex)
	if err != nil {
		return nil, err
	}

	if jErr != nil {
		result.Err.Message = jErr.Message
		result.Err.Code = jErr.Code
	}

	return &result, nil
}

func newTransaction() *transaction {
	return &transaction{}
}
