package security

import (
	"fmt"
	"os/exec"
	"strings"
)

// WindowsCredentialStore implements KeyStore using Windows Credential Manager
type WindowsCredentialStore struct{}

// newKeyStore creates a new WindowsCredentialStore for Windows
func newKeyStore() (KeyStore, error) {
	return &WindowsCredentialStore{}, nil
}

// NewWindowsCredentialStore creates a new WindowsCredentialStore
func NewWindowsCredentialStore() (*WindowsCredentialStore, error) {
	return &WindowsCredentialStore{}, nil
}

// StoreKey stores a key in Windows Credential Manager using cmdkey
func (w *WindowsCredentialStore) StoreKey(account string, key []byte) error {
	encoded := EncodeKey(key)
	targetName := fmt.Sprintf("%s/%s", KeystoreService, account)
	
	// cmdkey /generic:targetName /user:account /pass:encoded
	cmd := exec.Command("cmdkey",
		fmt.Sprintf("/generic:%s", targetName),
		fmt.Sprintf("/user:%s", account),
		fmt.Sprintf("/pass:%s", encoded))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to store key: %v, output: %s", err, string(output))
	}
	
	return nil
}

// RetrieveKey retrieves a key from Windows Credential Manager
// Note: This uses PowerShell to retrieve the credential as cmdkey doesn't support reading
func (w *WindowsCredentialStore) RetrieveKey(account string) ([]byte, error) {
	targetName := fmt.Sprintf("%s/%s", KeystoreService, account)
	
	// Use PowerShell to retrieve credential
	psCmd := fmt.Sprintf(`
$cred = Get-StoredCredential -Target "%s" -ErrorAction SilentlyContinue
if ($cred) { $cred.GetNetworkCredential().Password } else { exit 1 }
`, targetName)
	
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key: %v", err)
	}
	
	encoded := strings.TrimSpace(string(output))
	if encoded == "" {
		return nil, fmt.Errorf("key not found")
	}
	
	key, err := DecodeKey(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %v", err)
	}
	
	return key, nil
}

// DeleteKey removes a key from Windows Credential Manager
func (w *WindowsCredentialStore) DeleteKey(account string) error {
	targetName := fmt.Sprintf("%s/%s", KeystoreService, account)
	
	cmd := exec.Command("cmdkey",
		fmt.Sprintf("/delete:%s", targetName))
	
	_ = cmd.Run() // Ignore errors (key might not exist)
	return nil
}

// KeyExists checks if a key exists in Windows Credential Manager
func (w *WindowsCredentialStore) KeyExists(account string) bool {
	targetName := fmt.Sprintf("%s/%s", KeystoreService, account)
	
	psCmd := fmt.Sprintf(`
$cred = Get-StoredCredential -Target "%s" -ErrorAction SilentlyContinue
if ($cred) { exit 0 } else { exit 1 }
`, targetName)
	
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	err := cmd.Run()
	return err == nil
}

