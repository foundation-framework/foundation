package randutil

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
		panic("randutil: error occurred while generating random data: " + err.Error())
	}

	return u.String()
}

func Int(min, max int) int {
	if min > max {
		panic("randutil: min is greater than max")
	}

	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		panic("randutil: error occurred while generating random data: " + err.Error())
	}

	return int(binary.BigEndian.Uint32(bytes))%(max-min+1) + min
}

func IntRange(r [2]int) int {
	return Int(r[0], r[1])
}

func Hex(size int) string {
	return hex.EncodeToString(randomBytes(size))
}

func Base64(size int) string {
	return base64.StdEncoding.EncodeToString(randomBytes(size))
}

func randomBytes(size int) []byte {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic("randutil: error occurred while generating random data: " + err.Error())
	}

	return bytes
}
