// Package block defines our internal core domain logic related to Blocks.
// it should be exposes to external communication like tools, services (rest, graphql)

package block

import (
	"context"
	"github.com/INFURA/infra/internal"
	"github.com/INFURA/infra/repository"
)

func GetBlock(cxt context.Context, repository repository.Repository, blockNumber int64, hash string, details bool) (*internal.Block, error) {
	params := &internal.BqParams{
		Hash:        hash,
		BlockNumber: blockNumber,
		Details:     details,
	}
	return repository.GetBlock(cxt, params)
}
