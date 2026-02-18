package ethrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return c.rpc.SendTransaction(ctx, tx)
}

func (c *Client) WaitReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	deadline := time.Now().Add(2 * time.Minute)

	for time.Now().Before(deadline) {
		r, err := c.rpc.TransactionReceipt(ctx, txHash)
		if err == nil && r != nil {
			return r, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("timeout waiting for receipt: %s", txHash.Hex())
}
