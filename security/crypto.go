package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize   = 32 // AES-256
	nonceSize = 12 // GCM standard nonce size
	saltSize  = 32
	iterations = 100000
)

// MessageEncryptor handles encryption/decryption of messages between user and model
type MessageEncryptor struct {
	key []byte
}

// NewMessageEncryptor creates a new encryptor with the provided key
func NewMessageEncryptor(key []byte) (*MessageEncryptor, error) {
	if len(key) != keySize {
		return nil, fmt.Errorf("key must be %d bytes", keySize)
	}
	return &MessageEncryptor{key: key}, nil
}

// DeriveKeyFromPassword derives an encryption key from a password using PBKDF2
func DeriveKeyFromPassword(password string, salt []byte) []byte {
	if len(salt) != saltSize {
		panic("salt must be 32 bytes")
	}
	return pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)
}

// GenerateSalt generates a random salt for key derivation
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateKey generates a random encryption key
func GenerateKey() ([]byte, error) {
	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (e *MessageEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (e *MessageEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64-encoded ciphertext
func (e *MessageEncryptor) EncryptString(plaintext string) (string, error) {
	encrypted, err := e.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString decrypts a base64-encoded ciphertext string
func (e *MessageEncryptor) DecryptString(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decrypted, err := e.Decrypt(data)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

