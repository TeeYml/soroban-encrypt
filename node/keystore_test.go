package main

import (
	"os"
	"testing"
)

func TestKeyPersistence(t *testing.T) {
	dir := t.TempDir()

	priv1, err := LoadOrGenerateKey(dir)
	if err != nil {
		t.Fatalf("first load: %v", err)
	}

	priv2, err := LoadOrGenerateKey(dir)
	if err != nil {
		t.Fatalf("second load: %v", err)
	}

	if string(priv1.Bytes()) != string(priv2.Bytes()) {
		t.Error("key changed between restarts — persistence is broken")
	}

	if _, err := os.Stat(dir + "/node_key.pem"); err != nil {
		t.Errorf("PEM file not found after key generation: %v", err)
	}
}
