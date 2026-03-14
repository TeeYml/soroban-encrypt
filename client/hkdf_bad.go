package main

import (
	"crypto/sha256"
	"io"

	"golang.org/x/crypto/hkdf"
)

// WARNING: wrong info string — will cause decryption failures
const clientHKDFInfoWrong = "soroban-share-key-derivation"

func deriveAESKeyWrong(sharedSecret []byte) ([]byte, error) {
	r := hkdf.New(sha256.New, sharedSecret, nil, []byte(clientHKDFInfoWrong))
	key := make([]byte, 32)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}
	return key, nil
}
