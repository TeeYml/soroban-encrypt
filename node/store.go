package main

import (
	"encoding/json"
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

const (
	shareBucket = "shares"
	dataDir     = "./data"
)

var db *bolt.DB

// StoredShareDB extends StoredShare with metadata for persistence
type StoredShareDB struct {
	EphemeralPubKey []byte `json:"ephemeral_pubkey"`
	Ciphertext      []byte `json:"ciphertext"`
	Nonce           []byte `json:"nonce"`
}

func initDB() error {
	dir := os.Getenv("DATA_DIR")
	if dir == "" {
		dir = dataDir
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create data directory %s: %w", dir, err)
	}

	var err error
	db, err = bolt.Open(dir+"/shares.db", 0600, nil)
	if err != nil {
		return fmt.Errorf("failed to open BoltDB: %w", err)
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(shareBucket))
		if err != nil {
			return fmt.Errorf("failed to create shares bucket: %w", err)
		}
		return nil
	})
}

func saveShare(key string, share StoredShareDB) error {
	data, err := json.Marshal(share)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(shareBucket))
		return b.Put([]byte(key), data)
	})
}

func loadShare(key string) (StoredShareDB, bool, error) {
	var share StoredShareDB
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(shareBucket))
		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		return json.Unmarshal(v, &share)
	})
	if err != nil {
		return share, false, err
	}
	if share.Ciphertext == nil {
		return share, false, nil
	}
	return share, true, nil
}

func countShares() (int, error) {
	count := 0
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(shareBucket))
		return b.ForEach(func(k, v []byte) error {
			count++
			return nil
		})
	})
	return count, err
}

func deleteShare(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(shareBucket))
		return b.Delete([]byte(key))
	})
}
