package hashing

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/argon2"
)

const (
	saltLen = 8

	timeHash = 1
	memory   = 64 * 1024
	threads  = 4
	keyLen   = 32
)

func HashPass(plainPassword string) (string, error) {
	salt := make([]byte, saltLen)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("to rand read: %w", err)
	}

	return hex.EncodeToString(hashPassWithSalt(salt, plainPassword)), nil
}

func hashPassWithSalt(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, timeHash, memory, threads, keyLen)

	return append(salt, hashedPass...)
}

func ComparePassAndHash(passHashRaw string, plainPassword string) (bool, error) {
	passHash, err := hex.DecodeString(passHashRaw)
	if err != nil {
		return false, fmt.Errorf("to decode string: %w", err)
	}

	passHashCopy := make([]byte, len(passHash))

	copy(passHashCopy, passHash)

	salt := passHashCopy[0:saltLen]
	userPassHash := hashPassWithSalt(salt[:saltLen], plainPassword)

	return bytes.Equal(userPassHash, passHash), nil
}
