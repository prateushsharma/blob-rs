package ethrpc

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	rpc *ethclient.Client
}

func Dial(rpcURL string) (*Client, error) {
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: c}, nil
}

func (c *Client) Close() {
	c.rpc.Close()
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	return c.rpc.ChainID(ctx)
}

func (c *Client) PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	return c.rpc.PendingNonceAt(ctx, addr)
}

func (c *Client) WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, 15*time.Second)
}

func EnsureChainIDMatches(ctx context.Context, c *Client, expected *big.Int) error {
	got, err := c.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("rpc chainId: %w", err)
	}
	if got.Cmp(expected) != 0 {
		return fmt.Errorf("chainId mismatch: rpc=%s expected=%s", got.String(), expected.String())
	}
	return nil
}
