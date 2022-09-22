package infura

import (
	"context"
	model "github.com/INFURA/infra/internal"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type block struct {
}

func (b *block) getBlockByNumber(ctx context.Context, number int64, detail bool) (*model.Block, error) {
	client, err := getEthClient(ctx)
	if err != nil {
		return nil, err
	}
	txBlock, err := client.BlockByNumber(ctx, big.NewInt(number))
	if err != nil {
		return nil, err
	}

	result := &model.Block{
		Block: txBlock,
	}
	return result, nil
}

func (b *block) getBlockByHash(ctx context.Context, hash string, detail bool) (*model.Block, error) {
	client, err := getEthClient(ctx)
	if err != nil {
		return nil, err
	}
	txBlock, err := client.BlockByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	result := model.Block{
		Block: txBlock,
	}
	return &result, nil
}

func newBlock() *block {
	return &block{}
}
