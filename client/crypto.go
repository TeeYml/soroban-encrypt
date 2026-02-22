package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/stellar/go-stellar-sdk/keypair"
)

// GenerateAESKey generates a random 256-bit symmetric key.
func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	return key, nil
}

// SymmetricEncrypt runs AES-GCM encryption on input bytes.
func SymmetricEncrypt(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// SymmetricDecrypt decrypts AES-GCM ciphertext.
func SymmetricDecrypt(key []byte, ciphertext []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// ECIESEncrypt encrypts SSS shares under the node's public key.
func ECIESEncrypt(nodePubKeyBytes []byte, plaintext []byte) (ephemeralPubBytes []byte, ciphertext []byte, nonce []byte, err error) {
	// Parse the Node's P-256 public key
	nodePub, err := ecdh.P256().NewPublicKey(nodePubKeyBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid node public key: %w", err)
	}

	// Generate client ephemeral P-256 key pair
	ephemeralPriv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate ephemeral key: %w", err)
	}
	ephemeralPubBytes = ephemeralPriv.PublicKey().Bytes()

	// Compute shared secret using ECDH
	sharedSecret, err := ephemeralPriv.ECDH(nodePub)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("ECDH failed: %w", err)
	}

	// Encrypt the share using AES-256-GCM with the shared secret
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return nil, nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, err
	}

	nonce = make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, err
	}

	ciphertext = aesGCM.Seal(nil, nonce, plaintext, nil)
	return ephemeralPubBytes, ciphertext, nonce, nil
}

// SignPayload signs bytes using a private seed.
func SignPayload(stellarSeed string, payload []byte) (signature []byte, err error) {
	kp, err := keypair.ParseFull(stellarSeed)
	if err != nil {
		return nil, fmt.Errorf("invalid Stellar seed: %w", err)
	}

	sigBytes, err := kp.Sign(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}

	return sigBytes, nil
}

// GetAddressFromSeed gets the public G-address from a seed.
func GetAddressFromSeed(stellarSeed string) (string, error) {
	kp, err := keypair.ParseFull(stellarSeed)
	if err != nil {
		return "", err
	}
	return kp.Address(), nil
}
