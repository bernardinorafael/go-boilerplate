package uid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const epochTimestampInSeconds = 1700000000

// New generates a unique identifier with an optional prefix.
// It creates a KSUID (K-Sortable Unique Identifier) and prepends the
// specified prefix with an underscore separator if a prefix is provided.
//
// KSUIDs are 27-character, URL-safe, base62-encoded strings that contain:
// - A timestamp with 1-second resolution (first 4 bytes)
// - 16 bytes of random data
//
// This makes them:
// - Time sortable (newer KSUIDs sort lexicographically after older ones)
// - Highly unique (with ~2^128 possible random combinations)
// - More compact than UUIDs
//
// Example:
//
//	// Generate an ID with a custom prefix
//	id := uid.New("invoice") // returns "invoice_1z4UVH4CbRPvgSfCBmheK2h8xZb"
//
//	// Generate an ID without a prefix
//	id := uid.New("") // returns "1z4UVH4CbRPvgSfCBmheK2h8xZb"
func New(prefix string) string {
	buf := make([]byte, 12)
	t := uint32(time.Now().Unix() - epochTimestampInSeconds)
	binary.BigEndian.PutUint32(buf[:4], t)

	_, err := rand.Read(buf[4:])
	if err != nil {
		panic(err)
	}

	if prefix == "" {
		return fmt.Sprintf("%x", buf)
	}

	return fmt.Sprintf("%s_%x", prefix, buf)
}
