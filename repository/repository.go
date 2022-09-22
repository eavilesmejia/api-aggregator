// Repository package abstract the data source to our core domain, since
// the source of data should not be a concern for our business logic

package repository

import (
	"context"
	"github.com/INFURA/infra/internal"
)

type Repository interface {
	GetTransaction(ctx context.Context, params *internal.TxParams) (*internal.Transaction, error)
	GetBlock(ctx context.Context, params *internal.BqParams) (*internal.Block, error)
}
