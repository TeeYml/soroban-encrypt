package main

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestInitLoggerFallback(t *testing.T) {
	// Should not panic on invalid log level
	InitLogger("invalid-log-level")

	if zerolog.GlobalLevel() != zerolog.InfoLevel {
		t.Errorf("expected global level to fallback to InfoLevel, got %v", zerolog.GlobalLevel())
	}
}
