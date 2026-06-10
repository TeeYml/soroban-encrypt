package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSyncRequiresAuth(t *testing.T) {
	InitLogger("error")
	// /sync through apiKeyMiddleware with non-empty key
	handler := apiKeyMiddleware("testkey")(http.HandlerFunc(handleSync))
	req := httptest.NewRequest(http.MethodPost, "/sync", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without api key, got %d", rr.Code)
	}
}

func TestSyncLoopDetection(t *testing.T) {
	InitLogger("error")
	req := httptest.NewRequest(http.MethodPost, "/sync", nil)
	req.Header.Set("X-Forwarded-By", "peer-node-1")
	rr := httptest.NewRecorder()
	handleSync(rr, req)
	if rr.Code != http.StatusLoopDetected {
		t.Errorf("expected 508, got %d", rr.Code)
	}
}
