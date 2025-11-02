# Security Package

This package provides security features for SecLlama, including encryption, sandboxing, and secure key storage.

## Components

### Encryption (`crypto.go`)
- AES-256-GCM encryption for message content
- PBKDF2 key derivation from passwords
- Secure random key generation

### Key Storage (`keystore*.go`)
- Platform-specific secure key storage
- macOS: Keychain via `security` command
- Linux: Secret Service via `secret-tool`
- Windows: Credential Manager via `cmdkey`

### Sandboxing (`sandbox*.go`)
- Network isolation for runner processes
- Platform-specific implementations:
  - Linux: Network namespaces / unshare
  - macOS: sandbox-exec with custom profiles
  - Windows: Firewall rules

### Message Handling (`message_handler.go`)
- Encryption/decryption wrapper for messages
- JSON marshaling with encryption support

### HTTP Client (`http_client.go`)
- Localhost-only HTTP client
- Blocks all external connections
- Used for runner communication

### Security Manager (`manager.go`)
- Central security orchestration
- Key management lifecycle
- Sandbox configuration

## Usage

```go
import "github.com/ollama/ollama/security"

// Initialize security manager
mgr, err := security.GetManager()
if err != nil {
    log.Fatal(err)
}

// Encrypt a message
encrypted, err := mgr.EncryptMessage("secret data")

// Apply sandbox to command
cmd := exec.Command("runner", args...)
config := mgr.GetSandboxConfig(port)
err = security.ApplySandbox(cmd, config)

// Use localhost-only HTTP client
client := security.NewLocalhostOnlyClient()
```

## Testing

```bash
# Run tests
go test ./security/...

# Run with coverage
go test -cover ./security/...

# Run specific platform tests
go test -run TestSandbox ./security/
```

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| Encryption | ✅ | ✅ | ✅ |
| Key Storage | ✅ | ✅ | ✅ |
| Network Isolation | ✅ | ✅ | ⚠️ |
| Process Sandboxing | ✅ | ✅ | ⚠️ |

✅ Full support  
⚠️ Partial support (requires admin/manual setup)

## Security Considerations

1. **Key Storage**: Keys are stored in OS keychains which are protected by user authentication
2. **Network Isolation**: Requires appropriate OS permissions (may need root/admin)
3. **Encryption**: AES-256-GCM provides both confidentiality and authenticity
4. **Threat Model**: See [SECLLAMA_SECURITY.md](../SECLLAMA_SECURITY.md) for full threat model

## Contributing

When contributing security code:
1. Follow secure coding practices
2. Add comprehensive tests
3. Document security implications
4. Run security audit tools

## License

MIT License - see LICENSE file

