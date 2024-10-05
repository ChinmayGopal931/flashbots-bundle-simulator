package ethereum

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	ec *ethclient.Client
}

func NewClient(url string) (*Client, error) {
	ec, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{ec: ec}, nil
}

func (c *Client) GetBlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return c.ec.BlockByNumber(ctx, number)
}

func (c *Client) GetBalance(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.ec.BalanceAt(ctx, address, blockNumber)
}

func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.ec.SendTransaction(ctx, tx)
}

func (c *Client) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return c.ec.TransactionByHash(ctx, hash)
}

func (c *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.ec.TransactionReceipt(ctx, txHash)
}

func (c *Client) Close() {
	c.ec.Close()
}
