package randutil

import (
	"crypto/rand"
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

func Hex(size int) string {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic("randutil: error occurred while generating random data: " + err.Error())
	}

	return hex.EncodeToString(bytes)
}
