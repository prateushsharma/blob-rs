package tx

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
)

type BuildParams struct {
	ChainID *big.Int

	Nonce uint64
	To    common.Address
	Value *big.Int
	Data  []byte

	GasLimit uint64

	// EIP-1559
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int

	// EIP-4844 blob fee cap
	MaxFeePerBlobGas *big.Int

	// Blob artifacts
	Blob          kzg4844.Blob
	Commitment    kzg4844.Commitment
	Proof         kzg4844.Proof
	VersionedHash [32]byte
}

func BuildAndSignBlobTx(p BuildParams, privKey *ecdsa.PrivateKey) (*types.Transaction, common.Address, error) {
	// Defaults
	if p.Value == nil {
		p.Value = big.NewInt(0)
	}
	if p.Data == nil {
		p.Data = []byte{}
	}
	if p.GasLimit == 0 {
		p.GasLimit = 25000
	}

	// Convert big.Int -> uint256.Int (go-ethereum v1.17 uses uint256 for tx fields)
	chainIDU256, err := u256FromBig(p.ChainID)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid chainID: %w", err)
	}
	valueU256, err := u256FromBig(p.Value)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid value: %w", err)
	}
	feeCapU256, err := u256FromBig(p.MaxFeePerGas)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid maxFeePerGas: %w", err)
	}
	tipCapU256, err := u256FromBig(p.MaxPriorityFeePerGas)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid maxPriorityFeePerGas: %w", err)
	}
	blobFeeCapU256, err := u256FromBig(p.MaxFeePerBlobGas)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("invalid maxFeePerBlobGas: %w", err)
	}

	// Sidecar: carries blob data
	sidecar := &types.BlobTxSidecar{
		Blobs:       []kzg4844.Blob{p.Blob},
		Commitments: []kzg4844.Commitment{p.Commitment},
		Proofs:      []kzg4844.Proof{p.Proof},
	}

	// Versioned hash goes in tx body
	blobHashes := []common.Hash{
		common.BytesToHash(p.VersionedHash[:]),
	}

	txdata := &types.BlobTx{
		ChainID:    chainIDU256,
		Nonce:      p.Nonce,
		Gas:        p.GasLimit,
		To:         p.To,
		Value:      valueU256,
		Data:       p.Data,
		GasFeeCap:  feeCapU256,
		GasTipCap:  tipCapU256,
		BlobFeeCap: blobFeeCapU256,
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	}

	unsigned := types.NewTx(txdata)

	signer := types.LatestSignerForChainID(p.ChainID)
	signed, err := types.SignTx(unsigned, signer, privKey)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("sign tx: %w", err)
	}

	from, err := types.Sender(signer, signed)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("recover sender: %w", err)
	}

	return signed, from, nil
}

// u256FromBig converts big.Int to uint256.Int safely.
func u256FromBig(x *big.Int) (*uint256.Int, error) {
	if x == nil {
		return nil, fmt.Errorf("nil big.Int")
	}
	if x.Sign() < 0 {
		return nil, fmt.Errorf("negative value not allowed")
	}
	u, overflow := uint256.FromBig(x)
	if overflow {
		return nil, fmt.Errorf("uint256 overflow")
	}
	return u, nil
}
