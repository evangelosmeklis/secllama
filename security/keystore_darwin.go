package security

import (
	"fmt"
	"os/exec"
	"strings"
)

// KeychainKeyStore implements KeyStore using macOS Keychain
type KeychainKeyStore struct{}

// newKeyStore creates a new KeychainKeyStore for macOS
func newKeyStore() (KeyStore, error) {
	return &KeychainKeyStore{}, nil
}

// NewKeychainKeyStore creates a new KeychainKeyStore
func NewKeychainKeyStore() (*KeychainKeyStore, error) {
	return &KeychainKeyStore{}, nil
}

// StoreKey stores a key in macOS Keychain
func (k *KeychainKeyStore) StoreKey(account string, key []byte) error {
	encoded := EncodeKey(key)
	
	// First, try to delete any existing key
	_ = k.DeleteKey(account)
	
	// Add new key
	cmd := exec.Command("security", "add-generic-password",
		"-s", KeystoreService,
		"-a", account,
		"-w", encoded,
		"-U") // Update if exists
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store key in keychain: %v, output: %s", err, string(output))
	}
	
	return nil
}

// RetrieveKey retrieves a key from macOS Keychain
func (k *KeychainKeyStore) RetrieveKey(account string) ([]byte, error) {
	cmd := exec.Command("security", "find-generic-password",
		"-s", KeystoreService,
		"-a", account,
		"-w") // Output password only
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key from keychain: %v", err)
	}
	
	encoded := strings.TrimSpace(string(output))
	key, err := DecodeKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %v", err)
	}
	
	return key, nil
}

// DeleteKey removes a key from macOS Keychain
func (k *KeychainKeyStore) DeleteKey(account string) error {
	cmd := exec.Command("security", "delete-generic-password",
		"-s", KeystoreService,
		"-a", account)
	
	_ = cmd.Run() // Ignore errors (key might not exist)
	return nil
}

// KeyExists checks if a key exists in macOS Keychain
func (k *KeychainKeyStore) KeyExists(account string) bool {
	cmd := exec.Command("security", "find-generic-password",
		"-s", KeystoreService,
		"-a", account)
	
	err := cmd.Run()
	return err == nil
}

