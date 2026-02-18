package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/prateushsharma/blob-rs/internal/blob"
	"github.com/prateushsharma/blob-rs/internal/config"
	"github.com/prateushsharma/blob-rs/internal/ethrpc"
	"github.com/prateushsharma/blob-rs/internal/kzg"
	"github.com/prateushsharma/blob-rs/internal/tx"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "publish":
		runPublish(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func runPublish(args []string) {
	cfg, err := config.ParsePublishFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		usage()
		os.Exit(2)
	}

	payload, err := blob.ReadPayload(cfg.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}

	b, meta, err := blob.PackToBlob(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error packing blob: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("blob built ✅")
	fmt.Printf("payload=%d bytes (max=%d)\n", meta.PayloadLen, meta.MaxPayload)

	kzgres, err := kzg.Compute(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "kzg error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("kzg computed ✅")
	fmt.Printf("versioned_hash=%s\n", kzg.Hex32(kzgres.VersionedHash))

	rpcClient, err := ethrpc.Dial(cfg.RPCURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rpc dial error: %v\n", err)
		os.Exit(1)
	}
	defer rpcClient.Close()

	ctx, cancel := rpcClient.WithTimeout(context.Background())
	defer cancel()

	if err := ethrpc.EnsureChainIDMatches(ctx, rpcClient, cfg.ChainID); err != nil {
		fmt.Fprintf(os.Stderr, "chain error: %v\n", err)
		os.Exit(1)
	}

	priv := cfg.PrivateKey.(*ecdsa.PrivateKey)
	fromAddr := crypto.PubkeyToAddress(priv.PublicKey)

	nonce, err := rpcClient.PendingNonceAt(ctx, fromAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nonce error: %v\n", err)
		os.Exit(1)
	}

	signedTx, from, err := tx.BuildAndSignBlobTx(tx.BuildParams{
		ChainID:              cfg.ChainID,
		Nonce:                nonce,
		To:                   cfg.To,
		Value:                big.NewInt(0),
		Data:                 []byte{},
		GasLimit:             25000,
		MaxFeePerGas:         cfg.MaxFeePerGas,
		MaxPriorityFeePerGas: cfg.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     cfg.MaxFeePerBlobGas,
		Blob:                 b,
		Commitment:           kzgres.Commitment,
		Proof:                kzgres.Proof,
		VersionedHash:        kzgres.VersionedHash,
	}, priv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tx build error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("type3 tx signed ✅")
	fmt.Printf("from=%s\n", from.Hex())
	fmt.Printf("tx_hash=%s\n", signedTx.Hash().Hex())
	fmt.Printf("nonce=%d\n", nonce)

	// Broadcast
	if err := rpcClient.SendTransaction(ctx, signedTx); err != nil {
		fmt.Fprintf(os.Stderr, "send error: %v\n", err)
		fmt.Fprintln(os.Stderr, "NOTE: if this error mentions sidecar/blob, we may need a blob-specific send path.")
		os.Exit(1)
	}

	fmt.Println("tx sent ✅ waiting for receipt...")

	receipt, err := rpcClient.WaitReceipt(ctx, signedTx.Hash())
	if err != nil {
		fmt.Fprintf(os.Stderr, "receipt error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("mined ✅")
	fmt.Printf("status=%d block=%d\n", receipt.Status, receipt.BlockNumber.Uint64())
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  blobrs publish [flags]")
}
