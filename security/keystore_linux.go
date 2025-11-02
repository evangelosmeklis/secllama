package security

import (
	"fmt"
	"os/exec"
	"strings"
)

// SecretServiceKeyStore implements KeyStore using Linux Secret Service (gnome-keyring/kwallet)
type SecretServiceKeyStore struct{}

// newKeyStore creates a new SecretServiceKeyStore for Linux
func newKeyStore() (KeyStore, error) {
	// Check if secret-tool is available
	_, err := exec.LookPath("secret-tool")
	if err != nil {
		return nil, fmt.Errorf("secret-tool not found. Please install libsecret-tools package")
	}
	return &SecretServiceKeyStore{}, nil
}

// NewSecretServiceKeyStore creates a new SecretServiceKeyStore
func NewSecretServiceKeyStore() (*SecretServiceKeyStore, error) {
	// Check if secret-tool is available
	_, err := exec.LookPath("secret-tool")
	if err != nil {
		return nil, fmt.Errorf("secret-tool not found. Please install libsecret-tools package")
	}
	return &SecretServiceKeyStore{}, nil
}

// StoreKey stores a key using secret-tool
func (s *SecretServiceKeyStore) StoreKey(account string, key []byte) error {
	encoded := EncodeKey(key)
	
	cmd := exec.Command("secret-tool", "store",
		"--label", fmt.Sprintf("SecLlama %s", account),
		"service", KeystoreService,
		"account", account)
	
	cmd.Stdin = strings.NewReader(encoded)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store key: %v, output: %s", err, string(output))
	}
	
	return nil
}

// RetrieveKey retrieves a key using secret-tool
func (s *SecretServiceKeyStore) RetrieveKey(account string) ([]byte, error) {
	cmd := exec.Command("secret-tool", "lookup",
		"service", KeystoreService,
		"account", account)
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key: %v", err)
	}
	
	encoded := strings.TrimSpace(string(output))
	key, err := DecodeKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %v", err)
	}
	
	return key, nil
}

// DeleteKey removes a key using secret-tool
func (s *SecretServiceKeyStore) DeleteKey(account string) error {
	cmd := exec.Command("secret-tool", "clear",
		"service", KeystoreService,
		"account", account)
	
	_ = cmd.Run() // Ignore errors (key might not exist)
	return nil
}

// KeyExists checks if a key exists
func (s *SecretServiceKeyStore) KeyExists(account string) bool {
	cmd := exec.Command("secret-tool", "lookup",
		"service", KeystoreService,
		"account", account)
	
	err := cmd.Run()
	return err == nil
}

