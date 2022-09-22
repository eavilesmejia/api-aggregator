// Package infura defines and abstract all the methods to communicate to infura API.
// it includes serializers, deserializers, http request/response

package infura

import (
	"context"
	"github.com/INFURA/infra/internal"
)

type infura struct {
	transaction *transaction
	block       *block
}

func (i *infura) GetTransaction(ctx context.Context, params *internal.TxParams) (*internal.Transaction, error) {
	// check if we can use getTransactionByNumberAndIndex
	if (params.BlockString != "" || params.BlockNumber > 0) && params.TxIndex >= 0 {
		return i.transaction.getTransactionByNumberAndIndex(ctx, params.BlockString, params.BlockNumber, params.TxIndex)
	}

	if params.TxIndex >= 0 {
		return i.transaction.getTransactionByBlockHashAndIndex(ctx, params.Hash, params.TxIndex)
	}

	return i.transaction.getTransactionByHash(ctx, params.Hash)
}

func (i *infura) GetBlock(ctx context.Context, params *internal.BqParams) (*internal.Block, error) {
	if params.BlockNumber > 0 {
		return i.block.getBlockByNumber(ctx, params.BlockNumber, params.Details)
	}

	return i.block.getBlockByHash(ctx, params.Hash, params.Details)
}

func NewInfura() *infura {
	return &infura{transaction: newTransaction(), block: newBlock()}
}
