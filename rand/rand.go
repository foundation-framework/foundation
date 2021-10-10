package rand

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"

	"github.com/google/uuid"
)

func UUIDv1() string {
	u, err := uuid.NewUUID()
	if err != nil {
		panic("randutil: unexpected error: " + err.Error())
	}

	return u.String()
}

func UUIDv4() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic("rand: error occurred while generating random data: " + err.Error())
	}

	return u.String()
}

func Int(min, max int) int {
	if min > max {
		panic("rand: min is greater than max")
	}

	return int(binary.BigEndian.Uint32(Bytes(4)))%(max-min+1) + min
}

func IntBounds(r [2]int) int {
	return Int(r[0], r[1])
}

func Hex(size int) string {
	return hex.EncodeToString(Bytes(size))
}

func Base64(size int) string {
	return base64.StdEncoding.EncodeToString(Bytes(size))
}

func Bytes(size int) []byte {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic("rand: error occurred while generating random data: " + err.Error())
	}

	return bytes
}
