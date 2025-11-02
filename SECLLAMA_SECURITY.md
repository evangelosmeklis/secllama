# SecLlama Security Features

SecLlama is a security-enhanced fork of Ollama that implements multiple layers of protection to ensure private and secure AI model interactions.

## Overview

SecLlama provides three main security features:

1. **Network Isolation** - Models run in sandboxes with no internet access
2. **End-to-End Encryption** - All chat messages are encrypted
3. **Secure Key Storage** - Encryption keys are stored using OS-native keychains

## Security Features

### 1. Network Isolation & Sandboxing

All model runner processes are executed in isolated sandboxes that prevent any external network access.

**How it works:**
- **Linux**: Uses network namespaces (`CLONE_NEWNET`) or `unshare` command to isolate network
- **macOS**: Uses `sandbox-exec` with custom profiles blocking external network
- **Windows**: Uses Windows Firewall rules to block external connections

**Configuration:**
```bash
# Enable/disable sandboxing (default: enabled)
export SECLLAMA_ENABLE_SANDBOX=true

# Strict network isolation (default: enabled)
export SECLLAMA_STRICT_NETWORK_ISOLATION=true
```

**What's blocked:**
- All external network connections
- DNS resolution to external hosts
- Socket creation to non-localhost addresses

**What's allowed:**
- Localhost connections (required for runner communication)
- File system access to model files
- GPU access for acceleration

### 2. End-to-End Message Encryption

All chat messages between the user and model are encrypted using AES-256-GCM.

**How it works:**
- Each message is encrypted before transmission
- Uses AES-256 in GCM mode for authenticated encryption
- Keys are generated using cryptographically secure random number generator
- Automatic key management via OS keychain

**Configuration:**
```bash
# Enable/disable encryption (default: enabled)
export SECLLAMA_ENABLE_ENCRYPTION=true
```

**Encryption Details:**
- Algorithm: AES-256-GCM
- Key size: 256 bits
- Nonce size: 96 bits (12 bytes)
- Authentication tag: 128 bits

**Key Rotation:**
```bash
# Rotate encryption key (new messages use new key)
secllama key rotate
```

### 3. Secure Key Storage

Encryption keys are securely stored using OS-native credential management systems.

**Platform-specific storage:**
- **macOS**: Uses macOS Keychain via `security` command
- **Linux**: Uses Secret Service (gnome-keyring/kwallet) via `secret-tool`
- **Windows**: Uses Windows Credential Manager via `cmdkey`

**Key Management:**
```bash
# View key status
secllama key status

# Rotate key
secllama key rotate

# Delete keys (WARNING: will lose access to encrypted messages)
secllama key delete
```

**Key Properties:**
- Keys never stored in plain text on disk
- OS-level encryption and access control
- Per-user key isolation
- Automatic key generation on first run

## Security Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         User Application                      │
└──────────────────┬──────────────────────────────────────────┘
                   │ Encrypted Messages
                   │ (AES-256-GCM)
┌──────────────────▼──────────────────────────────────────────┐
│                    SecLlama Main Server                       │
│  ┌────────────────────────────────────────────────────────┐ │
│  │         Security Manager (Key Management)              │ │
│  └────────────────────────────────────────────────────────┘ │
│                        │                                      │
│                        │ Encrypted Messages                   │
│                        │ Localhost-only HTTP                  │
└────────────────────────┼──────────────────────────────────────┘
                         │
┌────────────────────────▼──────────────────────────────────────┐
│              Sandboxed Runner Process                          │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  Network Isolation:                                      │ │
│  │  ✓ No external network access                           │ │
│  │  ✓ Localhost-only communication                         │ │
│  │  ✓ Restricted file system access                        │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│                    ┌─────────────────┐                        │
│                    │   LLM Model     │                        │
│                    └─────────────────┘                        │
└────────────────────────────────────────────────────────────────┘
```

## Threat Model

SecLlama protects against:

1. **Data Exfiltration**: Models cannot send data to external servers
2. **Man-in-the-Middle Attacks**: Encrypted communication prevents eavesdropping
3. **Unauthorized Access**: Keys stored securely in OS keychain
4. **Model Tampering**: Sandboxing limits what models can do
5. **Network Reconnaissance**: No external DNS or network scanning

## Security Best Practices

### For Users

1. **Keep your system updated**: SecLlama relies on OS security features
2. **Use strong user account passwords**: Keychain access is protected by user auth
3. **Verify model sources**: Only use models from trusted sources
4. **Regular key rotation**: Rotate encryption keys periodically
5. **Monitor logs**: Check logs for security warnings

### For Developers

1. **Never disable security features in production**
2. **Audit model inputs**: Sanitize prompts before sending to models
3. **Implement rate limiting**: Prevent abuse of the API
4. **Use HTTPS**: When accessing SecLlama remotely, always use HTTPS
5. **Regular security audits**: Review security logs and configurations

## Limitations

### Current Limitations

1. **Multimodal content**: Images are not encrypted (base64 encoded only)
2. **Tool calls**: External tool calls are blocked by network isolation
3. **Remote models**: Remote model features are disabled for security
4. **Performance**: Encryption adds ~5-10ms latency per message
5. **Platform support**: Sandboxing quality varies by OS

### Known Trade-offs

- **Functionality vs Security**: Some Ollama features are disabled for security
- **Performance overhead**: Encryption and sandboxing add minor overhead
- **Compatibility**: Some third-party integrations may not work

## Troubleshooting

### Encryption Issues

**Problem**: "encryption key not found"
```bash
# Solution: Initialize new key
secllama key init
```

**Problem**: "failed to decrypt message"
```bash
# Solution: Key may have been rotated or corrupted
# Check key status
secllama key status
```

### Sandbox Issues

**Problem**: "failed to apply sandbox" (Linux)
```bash
# Solution: Install unshare tool
sudo apt-get install util-linux  # Debian/Ubuntu
sudo yum install util-linux      # RHEL/CentOS
```

**Problem**: "sandbox-exec not found" (macOS)
```bash
# Solution: sandbox-exec should be pre-installed on macOS
# Try with sudo or check macOS version
```

**Problem**: "Access denied" when creating firewall rules (Windows)
```bash
# Solution: Run as Administrator
# Right-click > Run as Administrator
```

### Network Isolation Issues

**Problem**: Model fails to start with network errors
```bash
# Solution: Disable strict isolation temporarily for debugging
export SECLLAMA_STRICT_NETWORK_ISOLATION=false
```

**Problem**: Runner can't communicate with main server
```bash
# Solution: Check firewall rules aren't blocking localhost
# Ensure port is not blocked
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SECLLAMA_ENABLE_ENCRYPTION` | `true` | Enable message encryption |
| `SECLLAMA_ENABLE_SANDBOX` | `true` | Enable process sandboxing |
| `SECLLAMA_STRICT_NETWORK_ISOLATION` | `true` | Enforce strict network isolation |

## Compliance & Certifications

SecLlama is designed with the following compliance frameworks in mind:

- **GDPR**: Data minimization and encryption at rest/in transit
- **HIPAA**: Encryption and access controls for healthcare data
- **SOC 2**: Security monitoring and access logging
- **ISO 27001**: Information security management

> **Note**: SecLlama provides security features but organizations must implement appropriate policies and procedures for full compliance.

## Security Reporting

If you discover a security vulnerability in SecLlama, please report it to:

- Email: security@secllama.example.com (replace with actual contact)
- GitHub Security Advisories: Use "Report a security vulnerability" on GitHub

Please do NOT create public issues for security vulnerabilities.

## License

SecLlama security features are licensed under MIT License. See LICENSE file for details.

## Contributing

We welcome security contributions! Please see CONTRIBUTING.md for guidelines.

For security-related contributions:
1. Review existing security architecture
2. Discuss major changes in an issue first
3. Include security test cases
4. Document security implications

## Acknowledgments

- Built on top of [Ollama](https://github.com/ollama/ollama)
- Uses industry-standard cryptography libraries
- Inspired by secure computing practices from various security frameworks

