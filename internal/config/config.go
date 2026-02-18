package config

import (
	"encoding/hex"
	"errors"
	"flag"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PublishConfig struct {
	RPCURL string

	ChainID *big.Int

	PrivateKeyHex string
	PrivateKey    any // we keep as any here; tx module will cast to *ecdsa.PrivateKey

	To   common.Address
	File string

	// EIP-1559 gas params (in wei)
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int

	// Blob fee param (in wei)
	MaxFeePerBlobGas *big.Int
}

func ParsePublishFlags(args []string) (*PublishConfig, error) {
	fs := flag.NewFlagSet("publish", flag.ContinueOnError)

	rpc := fs.String("rpc", "", "Execution RPC URL (e.g., https://sepolia.infura.io/v3/...)")
	chainID := fs.Int64("chain-id", 11155111, "Chain ID (Sepolia default: 11155111)")
	pk := fs.String("pk", "", "Hex private key (dev only; prefer env var in real use)")
	to := fs.String("to", "", "Recipient address (0x...)")
	file := fs.String("file", "", "Path to payload file to pack into a blob")

	maxFeeGwei := fs.Int64("max-fee-gwei", 30, "Max fee per gas in gwei (EIP-1559)")
	maxPrioGwei := fs.Int64("max-priority-fee-gwei", 2, "Max priority fee per gas in gwei (EIP-1559)")
	maxBlobFeeGwei := fs.Int64("max-blob-fee-gwei", 30, "Max fee per blob gas in gwei")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if strings.TrimSpace(*rpc) == "" {
		return nil, errors.New("missing --rpc")
	}
	if strings.TrimSpace(*pk) == "" {
		return nil, errors.New("missing --pk")
	}
	if strings.TrimSpace(*to) == "" {
		return nil, errors.New("missing --to")
	}
	if strings.TrimSpace(*file) == "" {
		return nil, errors.New("missing --file")
	}

	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(*pk, "0x"))
	if err != nil {
		return nil, errors.New("invalid --pk hex")
	}
	privKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, errors.New("invalid --pk (not a valid secp256k1 key)")
	}

	cfg := &PublishConfig{
		RPCURL:               *rpc,
		ChainID:              big.NewInt(*chainID),
		PrivateKeyHex:        *pk,
		PrivateKey:           privKey,
		To:                   common.HexToAddress(*to),
		File:                 *file,
		MaxFeePerGas:         gweiToWei(*maxFeeGwei),
		MaxPriorityFeePerGas: gweiToWei(*maxPrioGwei),
		MaxFeePerBlobGas:     gweiToWei(*maxBlobFeeGwei),
	}

	// Basic sanity check for address
	if cfg.To == (common.Address{}) {
		return nil, errors.New("invalid --to address")
	}

	return cfg, nil
}

func gweiToWei(g int64) *big.Int {
	// 1 gwei = 1e9 wei
	return new(big.Int).Mul(big.NewInt(g), big.NewInt(1_000_000_000))
}
