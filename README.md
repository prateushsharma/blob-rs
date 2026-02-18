# Blob-RS

> A minimal Go implementation for building, signing, and broadcasting EIP-4844 (Type-3) blob transactions on Ethereum.

Blob-RS is a lightweight infrastructure tool that handles the full blob transaction lifecycle — no SDK abstractions, no magic. Just raw KZG commitments, proper sidecar construction, and on-chain delivery.

---

## What It Does

- Packs arbitrary payloads into a 128 KiB EIP-4844 blob
- Computes KZG commitment and proof
- Derives the blob versioned hash
- Constructs and signs a Type-3 transaction with sidecar
- Broadcasts to an Ethereum network (e.g. Sepolia)
- Waits for and returns the transaction receipt

---

## ✨ Features — v0.1

| Feature | Status |
|---|---|
| Deterministic blob container format (`BR01`) | ✅ |
| Full KZG commitment + proof generation | ✅ |
| Versioned hash derivation | ✅ |
| Native `types.BlobTx` construction | ✅ |
| Sidecar attachment | ✅ |
| Transaction signing via go-ethereum | ✅ |
| Broadcast + receipt confirmation | ✅ |
| Sepolia (EIP-4844 enabled networks) | ✅ |

---

## 🏗 Architecture

```
blob-rs/
├── cmd/blobrs/         # CLI entry point
├── internal/
│   ├── config/         # Flag parsing + runtime config
│   ├── blob/           # Blob packing logic (BR01 container)
│   ├── kzg/            # Commitment + proof computation
│   ├── tx/             # Type-3 tx construction + signing
│   └── ethrpc/         # RPC client + broadcasting
```

---

## 📦 Blob Container Format — BR01

The blob uses a fixed **131,072-byte** layout:

| Offset | Size | Description |
|--------|------|-------------|
| `0` | `4` | Magic bytes: `BR01` |
| `4` | `8` | Payload length (uint64 LE) |
| `12` | `N` | Payload bytes |
| `12+N` | `...` | Zero padding to 128 KiB |

All unused bytes are zero-padded to exactly 128 KiB.

---

## 🚀 Quick Start (Sepolia)

### 1. Clone & Build

```bash
git clone https://github.com/<your-username>/blob-rs.git
cd blob-rs
go build ./...
```

### 2. Create a payload

```bash
echo "hello blob world" > payload.bin
```

### 3. Export your private key

```bash
export BLOBRS_PK_HEX="0xYOUR_PRIVATE_KEY"
```

> ⚠️ Never commit your private key. Always use environment variables.

### 4. Publish

```bash
go run ./cmd/blobrs publish \
  --rpc https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY \
  --chain-id 11155111 \
  --pk "$BLOBRS_PK_HEX" \
  --to 0x000000000000000000000000000000000000dEaD \
  --file ./payload.bin \
  --max-fee-gwei 40 \
  --max-priority-fee-gwei 2 \
  --max-blob-fee-gwei 40
```

**Expected output:**

```
blob built       ✅
kzg computed     ✅
type3 tx signed  ✅
tx sent          ✅
mined            ✅
status=1  block=...
```

---

## 🔎 Verifying the Transaction

| Tool | Link |
|------|------|
| Sepolia Etherscan | https://sepolia.etherscan.io |
| Blob Explorer | https://sepolia.blobscan.com |

Or query directly via RPC:

```bash
eth_getTransactionReceipt
```

---

## ⚙️ Requirements

- **Go** 1.24+
- **go-ethereum** v1.17+
- An RPC endpoint with EIP-4844 support
- Testnet ETH for execution gas + blob gas

---

## 📌 Important Notes

**Blob transactions require two separate fee markets:**
- Execution gas (standard EIP-1559)
- Blob gas (EIP-4844 blob fee market)

Blob data is **not** permanently stored in execution layer state — it is pruned after ~18 days. For permanent storage, consider a DA layer.

---

## 🧠 What This Demonstrates

Blob-RS is a reference implementation for understanding how EIP-4844 works at the protocol level:

- How KZG commitments integrate into Ethereum transactions
- How Type-3 transactions differ from legacy and EIP-1559 types
- How blob sidecars are attached and transmitted to nodes
- The full flow from raw bytes → on-chain blob hash

This makes it a foundational building block for:
- Blob aggregation systems
- Rollup data publishers
- Data availability tooling
- On-chain proof systems

---

## 🔮 Roadmap

### v0.2
- Multi-entry blob aggregation
- Merkle root inclusion
- Per-entry inclusion proofs
- JSON receipt output

### v0.3
- Blob fetch + on-chain verification
- Inclusion proof verification CLI
- Deterministic aggregation batching
- Advanced fee estimation

---

## 📜 License

[MIT](./LICENSE)
