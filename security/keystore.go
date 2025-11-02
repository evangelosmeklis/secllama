package security

import (
	"encoding/base64"
)

const (
	// KeystoreService is the service name used in keychain/credential manager
	KeystoreService = "secllama"
	// EncryptionKeyAccount is the account name for the main encryption key
	EncryptionKeyAccount = "message-encryption-key"
)

// KeyStore provides secure storage for encryption keys using OS-native mechanisms
type KeyStore interface {
	// StoreKey stores a key securely in the OS keychain/credential manager
	StoreKey(account string, key []byte) error
	// RetrieveKey retrieves a key from the OS keychain/credential manager
	RetrieveKey(account string) ([]byte, error)
	// DeleteKey removes a key from the OS keychain/credential manager
	DeleteKey(account string) error
	// KeyExists checks if a key exists in the keystore
	KeyExists(account string) bool
}

// GetKeyStore returns the appropriate KeyStore implementation for the current OS
// Implementation is in platform-specific files (keystore_*.go)
func GetKeyStore() (KeyStore, error) {
	return newKeyStore()
}

// Base64Key helper functions for encoding/decoding keys
func EncodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func DecodeKey(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

