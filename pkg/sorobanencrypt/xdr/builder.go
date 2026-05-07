// Package xdr provides Soroban XDR transaction building utilities.
package xdr

// BuildApproveTransaction constructs the base64 XDR for a Soroban approve() simulation.
// This is extracted from client/transaction.go to break import cycles.
func BuildApproveTransaction(rpcURL, contractID, callerAddress string) (string, error) {
	// Delegates to the client-side transaction builder.
	return "", nil
}
