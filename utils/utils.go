package utils

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/google/uuid"
)

func GetUUID() string {
	clientID := uuid.New().String()
	return clientID
}

func GetUint64UUID() uint64 {
	var buf [8]byte
	_, _ = rand.Read(buf[:])
	return binary.BigEndian.Uint64(buf[:])
}
