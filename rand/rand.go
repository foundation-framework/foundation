package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"

	"github.com/foundation-framework/foundation/errors"
	"github.com/google/uuid"
)

const (
	SizeUUID = 16
)

func Bytes(size int) []byte {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		errors.Panicf("rand: unexpected generation error: %s", err)
	}

	return bytes
}

func UUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		errors.Panicf("rand: unexpected uuid error: %s", err)
	}

	return u.String()
}

func Int(min, max int) int {
	if min > max {
		min, max = max, min
	}

	return int(binary.BigEndian.Uint32(Bytes(4)))%(max-min+1) + min
}

func Hex(size int) string {
	return hex.EncodeToString(Bytes(size))
}

func Base64(size int) string {
	return base64.StdEncoding.EncodeToString(Bytes(size))
}
