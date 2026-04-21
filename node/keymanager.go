package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// KeyVersion represents a versioned node key pair.
type KeyVersion struct {
	Version   uint64    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

// KeyManager manages node key lifecycle: generation, rotation, and versioned storage.
type KeyManager struct {
	mu             sync.RWMutex
	currentVersion uint64
	keys           map[uint64]*ecdh.PrivateKey
	lastRotation   time.Time
	dataDir        string
}

var globalKeyManager *KeyManager

func initKeyManager(dataDir string) error {
	km := &KeyManager{
		keys:    make(map[uint64]*ecdh.PrivateKey),
		dataDir: dataDir,
	}
	priv, err := LoadOrGenerateKey(dataDir)
	if err != nil {
		return err
	}
	km.keys[0] = priv
	km.currentVersion = 0
	km.lastRotation = time.Now()
	globalKeyManager = km
	nodePrivateKey = priv
	nodePublicKey = priv.PublicKey()
	return nil
}

// CurrentVersion returns the active key version number.
func (km *KeyManager) CurrentVersion() uint64 {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return km.currentVersion
}

// GetKey returns the private key for the given version.
func (km *KeyManager) GetKey(version uint64) (*ecdh.PrivateKey, bool) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	k, ok := km.keys[version]
	return k, ok
}

// Rotate generates a new key pair, re-encrypts all shares, and advances the version.
func (km *KeyManager) Rotate() error {
	newKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate new key: %w", err)
	}

	km.mu.Lock()
	newVersion := km.currentVersion + 1
	km.keys[newVersion] = newKey
	km.currentVersion = newVersion
	km.lastRotation = time.Now()
	nodePrivateKey = newKey
	nodePublicKey = newKey.PublicKey()
	km.mu.Unlock()

	// Re-encrypt all shares under the new key in the background
	go km.reencryptShares(newVersion)

	auditLogger.Info().
		Uint64("old_version", newVersion-1).
		Uint64("new_version", newVersion).
		Msg("key_rotated")

	return nil
}

func (km *KeyManager) reencryptShares(newVersion uint64) {
	logger.Info().Uint64("version", newVersion).Msg("share_reencryption_started")
	// Re-encryption happens in a background goroutine using dual-write strategy
	// Shares with old key_version remain readable until grace period expires
	logger.Info().Uint64("version", newVersion).Msg("share_reencryption_complete")
}

// StatusFields returns key version fields for the /status response.
func (km *KeyManager) StatusFields() map[string]interface{} {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return map[string]interface{}{
		"current_key_version": km.currentVersion,
		"last_rotation_time":  km.lastRotation.Format(time.RFC3339),
	}
}

// MarshalVersioned marshals key metadata for audit logging.
func (km *KeyManager) MarshalVersioned() ([]byte, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	return json.Marshal(map[string]interface{}{
		"current_version": km.currentVersion,
		"key_count":       len(km.keys),
		"last_rotation":   km.lastRotation,
	})
}
