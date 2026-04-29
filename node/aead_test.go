package main

import (
	"bytes"
	"testing"
)

func TestBuildAAD(t *testing.T) {
	aad := BuildAAD("CONTRACT1", "OBJECT1")
	expected := []byte("soroban-encrypt|CONTRACT1|OBJECT1")
	if !bytes.Equal(aad, expected) {
		t.Errorf("AAD mismatch: got %q, want %q", aad, expected)
	}
}

func TestBuildAADUniqueness(t *testing.T) {
	aad1 := BuildAAD("C1", "O1")
	aad2 := BuildAAD("C1", "O2")
	aad3 := BuildAAD("C2", "O1")
	if bytes.Equal(aad1, aad2) || bytes.Equal(aad1, aad3) {
		t.Error("AAD must differ for different (contract, object) pairs")
	}
}
