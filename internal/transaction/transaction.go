// Package transaction defines our internal core domain logic related to Transactions.
// it should be exposes to external communication like tools, services (rest, graphql)

package transaction

import (
	"context"
	"github.com/INFURA/infra/internal"
	"github.com/INFURA/infra/repository"
)

func GetTransaction(ctx context.Context, repository repository.Repository,
	blockString string, blockNumber int64, txIndex int64, txHash string) (*internal.Transaction, error) {
	params := &internal.TxParams{
		BlockString: blockString,
		BlockNumber: blockNumber,
		TxIndex:     txIndex,
		Hash:        txHash,
	}
	return repository.GetTransaction(ctx, params)
}
