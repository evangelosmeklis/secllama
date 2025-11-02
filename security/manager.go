package security

import (
	"fmt"
	"log/slog"
	"sync"
)

// Manager handles security operations for secllama
type Manager struct {
	keyStore  KeyStore
	encryptor *MessageEncryptor
	mu        sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton security manager
func GetManager() (*Manager, error) {
	var err error
	once.Do(func() {
		instance, err = newManager()
	})
	return instance, err
}

// newManager creates a new security manager
func newManager() (*Manager, error) {
	keyStore, err := GetKeyStore()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keystore: %v", err)
	}
	
	m := &Manager{
		keyStore: keyStore,
	}
	
	// Try to load existing encryption key, or create a new one
	err = m.initializeEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption key: %v", err)
	}
	
	return m, nil
}

// initializeEncryptionKey loads or creates the encryption key
func (m *Manager) initializeEncryptionKey() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var key []byte
	var err error
	
	// Try to retrieve existing key
	if m.keyStore.KeyExists(EncryptionKeyAccount) {
		key, err = m.keyStore.RetrieveKey(EncryptionKeyAccount)
		if err != nil {
			slog.Warn("failed to retrieve existing key, creating new one", "error", err)
		} else {
			slog.Info("loaded existing encryption key from keystore")
		}
	}
	
	// Create new key if needed
	if key == nil {
		key, err = GenerateKey()
		if err != nil {
			return fmt.Errorf("failed to generate encryption key: %v", err)
		}
		
		err = m.keyStore.StoreKey(EncryptionKeyAccount, key)
		if err != nil {
			return fmt.Errorf("failed to store encryption key: %v", err)
		}
		
		slog.Info("generated and stored new encryption key")
	}
	
	// Create encryptor
	m.encryptor, err = NewMessageEncryptor(key)
	if err != nil {
		return fmt.Errorf("failed to create message encryptor: %v", err)
	}
	
	return nil
}

// EncryptMessage encrypts a message
func (m *Manager) EncryptMessage(plaintext string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.encryptor == nil {
		return "", fmt.Errorf("encryptor not initialized")
	}
	
	return m.encryptor.EncryptString(plaintext)
}

// DecryptMessage decrypts a message
func (m *Manager) DecryptMessage(ciphertext string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.encryptor == nil {
		return "", fmt.Errorf("encryptor not initialized")
	}
	
	return m.encryptor.DecryptString(ciphertext)
}

// RotateKey generates a new encryption key and re-encrypts data
func (m *Manager) RotateKey() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Generate new key
	newKey, err := GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate new key: %v", err)
	}
	
	// Store new key
	err = m.keyStore.StoreKey(EncryptionKeyAccount, newKey)
	if err != nil {
		return fmt.Errorf("failed to store new key: %v", err)
	}
	
	// Create new encryptor
	m.encryptor, err = NewMessageEncryptor(newKey)
	if err != nil {
		return fmt.Errorf("failed to create new encryptor: %v", err)
	}
	
	slog.Info("encryption key rotated successfully")
	return nil
}

// GetSandboxConfig returns the default sandbox configuration
func (m *Manager) GetSandboxConfig(port int) SandboxConfig {
	return SandboxConfig{
		AllowLocalhost:   true,
		AllowedPorts:     []int{port},
		WorkingDirectory: "/tmp/secllama",
		AllowedReadPaths: []string{
			"/tmp/secllama",
			"/var/tmp/secllama",
		},
		AllowedWritePaths: []string{
			"/tmp/secllama",
			"/var/tmp/secllama",
		},
	}
}

// IsEncryptionEnabled returns whether encryption is enabled
func (m *Manager) IsEncryptionEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.encryptor != nil
}

