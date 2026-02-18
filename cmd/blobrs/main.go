package main

import (
	"fmt"
	"os"

	"blob-rs/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	switch cmd {
	case "publish":
		cfg, err := config.ParsePublishFlags(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
			publishUsage()
			os.Exit(2)
		}

		// v0.1: wiring only (tx publishing comes in next commits)
		fmt.Println("publish config loaded ✅")
		fmt.Printf("rpc=%s chainID=%s to=%s file=%s\n", cfg.RPCURL, cfg.ChainID.String(), cfg.To.Hex(), cfg.File)

	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  blobrs publish [flags]")
	fmt.Println("")
	publishUsage()
}

func publishUsage() {
	fmt.Println("Publish flags:")
	fmt.Println("  --rpc <url>                     Execution RPC URL")
	fmt.Println("  --chain-id <id>                 Chain ID (default: 11155111 Sepolia)")
	fmt.Println("  --pk <hex>                      Private key hex (dev only)")
	fmt.Println("  --to <0xaddr>                   Recipient address")
	fmt.Println("  --file <path>                   Payload file to pack into a blob")
	fmt.Println("  --max-fee-gwei <n>              Max fee per gas (default 30)")
	fmt.Println("  --max-priority-fee-gwei <n>     Max priority fee per gas (default 2)")
	fmt.Println("  --max-blob-fee-gwei <n>         Max fee per blob gas (default 30)")
}
