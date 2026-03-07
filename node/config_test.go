package main

import (
	"os"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{"defaults are valid", nil, false},
		{"auto tls without domain fails", map[string]string{"TLS_MODE": "auto"}, true},
		{"auto tls with domain ok", map[string]string{"TLS_MODE": "auto", "DOMAIN": "example.com"}, false},
		{"manual tls no certs fails", map[string]string{"TLS_MODE": "manual"}, true},
		{"manual tls with certs ok", map[string]string{"TLS_MODE": "manual", "TLS_CERT": "c.pem", "TLS_KEY": "k.pem"}, false},
		{"invalid tls mode fails", map[string]string{"TLS_MODE": "bogus"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}
			err := LoadConfig().Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() err=%v wantErr=%v", err, tt.wantErr)
			}
		})
	}
}
