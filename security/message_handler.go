package security

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// SecureMessage wraps a message with encryption metadata
type SecureMessage struct {
	Content   string `json:"content"`
	Encrypted bool   `json:"encrypted"`
}

// EncryptMessageContent encrypts message content if encryption is enabled
func EncryptMessageContent(content string, enableEncryption bool) (string, bool, error) {
	if !enableEncryption {
		return content, false, nil
	}
	
	mgr, err := GetManager()
	if err != nil {
		slog.Warn("encryption enabled but manager not available", "error", err)
		return content, false, nil
	}
	
	encrypted, err := mgr.EncryptMessage(content)
	if err != nil {
		return "", false, fmt.Errorf("failed to encrypt message: %v", err)
	}
	
	return encrypted, true, nil
}

// DecryptMessageContent decrypts message content if it's encrypted
func DecryptMessageContent(content string, isEncrypted bool) (string, error) {
	if !isEncrypted {
		return content, nil
	}
	
	mgr, err := GetManager()
	if err != nil {
		return "", fmt.Errorf("decryption required but manager not available: %v", err)
	}
	
	decrypted, err := mgr.DecryptMessage(content)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}
	
	return decrypted, nil
}

// SecureJSONMarshaler provides encryption-aware JSON marshaling
type SecureJSONMarshaler struct {
	enableEncryption bool
}

// NewSecureJSONMarshaler creates a new secure JSON marshaler
func NewSecureJSONMarshaler(enableEncryption bool) *SecureJSONMarshaler {
	return &SecureJSONMarshaler{enableEncryption: enableEncryption}
}

// MarshalSecureMessage marshals a message with optional encryption
func (m *SecureJSONMarshaler) MarshalSecureMessage(msg interface{}) ([]byte, error) {
	// If encryption is disabled, just marshal normally
	if !m.enableEncryption {
		return json.Marshal(msg)
	}
	
	// Marshal to JSON first
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	
	// Encrypt the JSON data
	mgr, err := GetManager()
	if err != nil {
		slog.Warn("encryption enabled but manager not available", "error", err)
		return data, nil
	}
	
	encrypted, err := mgr.EncryptMessage(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %v", err)
	}
	
	// Create a secure envelope
	envelope := SecureMessage{
		Content:   encrypted,
		Encrypted: true,
	}
	
	return json.Marshal(envelope)
}

// UnmarshalSecureMessage unmarshals a message with optional decryption
func (m *SecureJSONMarshaler) UnmarshalSecureMessage(data []byte, v interface{}) error {
	// Try to unmarshal as SecureMessage first
	var envelope SecureMessage
	if err := json.Unmarshal(data, &envelope); err == nil && envelope.Encrypted {
		// This is an encrypted message, decrypt it
		mgr, err := GetManager()
		if err != nil {
			return fmt.Errorf("decryption required but manager not available: %v", err)
		}
		
		decrypted, err := mgr.DecryptMessage(envelope.Content)
		if err != nil {
			return fmt.Errorf("failed to decrypt message: %v", err)
		}
		
		// Unmarshal the decrypted content
		return json.Unmarshal([]byte(decrypted), v)
	}
	
	// Not encrypted, unmarshal normally
	return json.Unmarshal(data, v)
}

