package uid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const epochTimestampInSeconds = 1700000000

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
