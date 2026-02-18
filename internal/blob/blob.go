package blob

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	// EIP-4844 blob size is fixed: 131072 bytes (128 KiB).
	BlobSize = 131072

	// v0.1 container header:
	// [0..3]  magic "BR01"
	// [4..11] payload length (uint64 little-endian)
	HeaderSize = 12

	Magic0 = 'B'
	Magic1 = 'R'
	Magic2 = '0'
	Magic3 = '1'
)

type Meta struct {
	PayloadLen uint64
	MaxPayload uint64
}

// ReadPayload reads the raw bytes from the file at path.
func ReadPayload(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// PackToBlob packs payload into a fixed-size EIP-4844 blob.
// Layout:
//   - 4 bytes magic "BR01"
//   - 8 bytes payload length (uint64 LE)
//   - payload bytes
//   - zero padding to 131072 bytes
func PackToBlob(payload []byte) (kzg4844.Blob, Meta, error) {
	var out kzg4844.Blob

	maxPayload := BlobSize - HeaderSize
	if len(payload) > maxPayload {
		return out, Meta{}, fmt.Errorf("payload too large: %d bytes (max %d)", len(payload), maxPayload)
	}

	// Header: magic
	out[0] = Magic0
	out[1] = Magic1
	out[2] = Magic2
	out[3] = Magic3

	// Header: payload length
	binary.LittleEndian.PutUint64(out[4:12], uint64(len(payload)))

	// Body: payload
	copy(out[HeaderSize:HeaderSize+len(payload)], payload)

	// Remaining bytes are already zero (Go default) => padding.
	return out, Meta{
		PayloadLen: uint64(len(payload)),
		MaxPayload: uint64(maxPayload),
	}, nil
}

// UnpackFromBlob is optional but handy for local testing and later receipts.
// It validates the header and extracts payload.
func UnpackFromBlob(b kzg4844.Blob) ([]byte, Meta, error) {
	if b[0] != Magic0 || b[1] != Magic1 || b[2] != Magic2 || b[3] != Magic3 {
		return nil, Meta{}, fmt.Errorf("invalid blob magic header")
	}
	n := binary.LittleEndian.Uint64(b[4:12])

	maxPayload := uint64(BlobSize - HeaderSize)
	if n > maxPayload {
		return nil, Meta{}, fmt.Errorf("invalid payload length in blob: %d (max %d)", n, maxPayload)
	}
	payload := make([]byte, n)
	copy(payload, b[HeaderSize:HeaderSize+n])

	return payload, Meta{PayloadLen: n, MaxPayload: maxPayload}, nil
}
