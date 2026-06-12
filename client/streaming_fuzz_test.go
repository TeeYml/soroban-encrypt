//go:build go1.18

package main

import (
	"crypto/rand"
	"os"
	"testing"
)

func FuzzDecryptFile(f *testing.F) {
	// Add seed corpus entries
	f.Add([]byte("SE01\x00\x00\x00\x00"))
	f.Add([]byte("SE02\x01\x00\x00\x00\x10\x00\x00\x00"))
	f.Add([]byte("BADD"))
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		key := make([]byte, 32)
		rand.Read(key)
		tmp := t.TempDir()
		in := tmp + "/in.enc"
		out := tmp + "/out.bin"
		os.WriteFile(in, data, 0600)
		// Must not panic regardless of input
		_ = DecryptFileStreaming(in, out, key)
	})
}
