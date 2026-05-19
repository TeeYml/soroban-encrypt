package main

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"
)

func makeKey(t *testing.T) []byte {
	t.Helper()
	k := make([]byte, 32)
	rand.Read(k)
	return k
}

func TestStreamingRoundTrip(t *testing.T) {
	key := makeKey(t)
	plain := make([]byte, 5*1024*1024) // 5MB
	rand.Read(plain)

	tmp := t.TempDir()
	src := tmp + "/plain.bin"
	enc := tmp + "/plain.bin.enc"
	dec := tmp + "/plain.bin.dec"

	os.WriteFile(src, plain, 0600)

	if err := EncryptFileStreaming(src, enc, key, false); err != nil {
		t.Fatalf("EncryptFileStreaming: %v", err)
	}
	if err := DecryptFileStreaming(enc, dec, key); err != nil {
		t.Fatalf("DecryptFileStreaming: %v", err)
	}

	got, _ := os.ReadFile(dec)
	if !bytes.Equal(got, plain) {
		t.Error("decrypted output does not match original plaintext")
	}
}

func TestDecryptTruncated(t *testing.T) {
	key := makeKey(t)
	tmp := t.TempDir()
	// Write a file with v2 magic but truncated body
	f := tmp + "/truncated.enc"
	os.WriteFile(f, []byte("SE02\x05\x00\x00\x00"), 0600)
	if err := DecryptFileStreaming(f, tmp+"/out.bin", key); err == nil {
		t.Error("expected error on truncated .enc file, got nil")
	}
}

func TestDecryptUnknownMagic(t *testing.T) {
	key := makeKey(t)
	tmp := t.TempDir()
	f := tmp + "/bad.enc"
	os.WriteFile(f, []byte("BADD\x00\x00\x00\x00"), 0600)
	if err := DecryptFileStreaming(f, tmp+"/out.bin", key); err == nil {
		t.Error("expected error on unknown magic, got nil")
	}
}
