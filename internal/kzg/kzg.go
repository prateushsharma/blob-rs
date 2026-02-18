package kzg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

// Version byte for EIP-4844 blob versioned hashes (current is 0x01).
const VersionedHashVersionByte = 0x01

type Result struct {
	Commitment kzg4844.Commitment
	Proof      kzg4844.Proof
	// 32 bytes: versioned hash as used in the type-3 tx body
	VersionedHash [32]byte
}

// Compute computes commitment + proof for a blob, then derives the versioned hash.
// Also runs a local Verify() when possible (nice sanity check).
func Compute(blob kzg4844.Blob) (Result, error) {
	commitment, err := kzg4844.BlobToCommitment(&blob)
	if err != nil {
		return Result{}, fmt.Errorf("BlobToCommitment: %w", err)
	}

	proof, err := kzg4844.ComputeBlobProof(&blob, commitment)
	if err != nil {
		return Result{}, fmt.Errorf("ComputeBlobProof: %w", err)
	}

	// Optional sanity check
	if err := kzg4844.VerifyBlobProof(&blob, commitment, proof); err != nil {
		return Result{}, fmt.Errorf("VerifyBlobProof failed: %w", err)
	}

	vh := commitmentToVersionedHash(commitment)

	return Result{
		Commitment:    commitment,
		Proof:         proof,
		VersionedHash: vh,
	}, nil
}

// EIP-4844 versioned hash derivation:
// versioned_hash = version_byte || sha256(commitment)[1..31]
func commitmentToVersionedHash(c kzg4844.Commitment) [32]byte {
	h := sha256.Sum256(c[:])
	var out [32]byte
	out[0] = VersionedHashVersionByte
	copy(out[1:], h[1:])
	return out
}

func Hex32(b [32]byte) string { return "0x" + hex.EncodeToString(b[:]) }
func Hex48(b [48]byte) string { return "0x" + hex.EncodeToString(b[:]) }
