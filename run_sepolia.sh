#!/usr/bin/env bash
set -euo pipefail

# Required:
: "${BLOBRS_PK_HEX:?Set BLOBRS_PK_HEX to your private key hex (0x...)}"

# Optional (default: your Alchemy Sepolia RPC)
BLOBRS_RPC_URL="${BLOBRS_RPC_URL:-https://eth-sepolia.g.alchemy.com/v2/AIdUeGRr3tnVrzfea-WpKoEEBQ1SZDfy}"

# Payload file (default: ./payload.bin)
PAYLOAD="${1:-./payload.bin}"

# Quick payload if missing
if [ ! -f "$PAYLOAD" ]; then
  echo "hello blob world" > "$PAYLOAD"
fi

go run ./cmd/blobrs publish \
  --rpc "$BLOBRS_RPC_URL" \
  --chain-id 11155111 \
  --pk "$BLOBRS_PK_HEX" \
  --to 0x000000000000000000000000000000000000dEaD \
  --file "$PAYLOAD" \
  --max-fee-gwei 40 \
  --max-priority-fee-gwei 2 \
  --max-blob-fee-gwei 40
