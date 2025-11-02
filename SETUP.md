# SecLlama Setup Guide

This guide will help you get SecLlama up and running with all security features enabled.

## Prerequisites

### All Platforms
- Go 1.22 or later
- Git

### Linux
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install build-essential util-linux libsecret-tools

# RHEL/Fedora/CentOS
sudo yum install gcc gcc-c++ make util-linux libsecret
```

### macOS
```bash
# Install Xcode Command Line Tools
xcode-select --install

# sandbox-exec is pre-installed on macOS
```

### Windows
```bash
# Install Visual Studio Build Tools or MinGW
# Windows Credential Manager is pre-installed
```

## Installation

### From Source

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/secllama.git
cd secllama
```

2. **Build SecLlama**
```bash
# Build the project
go generate ./...
go build .

# This will create the 'secllama' executable (or 'secllama.exe' on Windows)
```

3. **Install (optional)**
```bash
# Linux/macOS
sudo mv secllama /usr/local/bin/

# Windows - add to PATH or move to desired location
```

## Initial Setup

### 1. Initialize Security Keys

On first run, SecLlama will automatically:
- Generate encryption keys
- Store them securely in your OS keychain
- Apply default security policies

```bash
# Start SecLlama server
./secllama serve

# You should see:
# "Generated and stored new encryption key"
# "Security features initialized successfully"
```

### 2. Verify Security Status

```bash
# Check security configuration
./secllama security status

# Expected output:
# ✓ Encryption: Enabled
# ✓ Sandboxing: Enabled  
# ✓ Network Isolation: Enabled
# ✓ Key Storage: OS Keychain
```

### 3. Test the Installation

```bash
# Pull a test model (this download happens over normal internet)
./secllama pull llama3.2:1b

# Run a test chat (this runs in isolated sandbox)
./secllama run llama3.2:1b "Hello, are you secure?"

# You should see encrypted communication in debug logs
```

## Configuration

### Environment Variables

Create a `.secllama` configuration file in your home directory:

```bash
# ~/.secllama/config

# Security settings (all default to true)
SECLLAMA_ENABLE_ENCRYPTION=true
SECLLAMA_ENABLE_SANDBOX=true
SECLLAMA_STRICT_NETWORK_ISOLATION=true

# Ollama compatibility settings
OLLAMA_HOST=127.0.0.1:11434
OLLAMA_MODELS=/path/to/models
```

### Platform-Specific Configuration

#### Linux

For best security on Linux, ensure network namespaces are available:

```bash
# Check if namespaces are supported
unshare --help | grep "net"

# If not available, install util-linux
sudo apt-get install util-linux

# For non-root users, you may need to enable user namespaces
echo "kernel.unprivileged_userns_clone=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

#### macOS

macOS sandboxing should work out of the box. Verify:

```bash
# Check if sandbox-exec is available
which sandbox-exec

# Should output: /usr/bin/sandbox-exec
```

#### Windows

For Windows Firewall integration (requires Administrator):

```bash
# Run PowerShell as Administrator
# The firewall rules will be created automatically when you start SecLlama
```

## Verifying Security

### 1. Test Network Isolation

```bash
# Start SecLlama with debug logging
SECLLAMA_LOG_LEVEL=debug ./secllama serve

# In another terminal, try to use a model
./secllama run llama3.2:1b "Make a request to httpbin.org"

# The model should fail to make any external connections
# You should see in logs: "connection blocked by security policy"
```

### 2. Test Encryption

```bash
# Enable debug logging to see encryption in action
SECLLAMA_LOG_LEVEL=debug ./secllama run llama3.2:1b "Test message"

# In the logs, look for:
# "encrypted message: <base64-encoded-ciphertext>"
# "decrypted message: Test message"
```

### 3. Test Key Storage

```bash
# Check keychain (macOS)
security find-generic-password -s secllama -a message-encryption-key

# Check keychain (Linux)
secret-tool search service secllama account message-encryption-key

# Check credential manager (Windows PowerShell)
Get-StoredCredential -Target "secllama/message-encryption-key"
```

## Usage Examples

### Basic Chat

```bash
# Start a chat session
./secllama run llama3.2

>>> Hello!
# Messages are automatically encrypted
```

### API Usage

```bash
# Start server
./secllama serve

# Use API (messages are encrypted)
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [{"role": "user", "content": "Hello"}]
}'
```

### Python Client

```python
import ollama

# SecLlama is API-compatible with Ollama
client = ollama.Client(host='http://localhost:11434')

# Messages are automatically encrypted by SecLlama server
response = client.chat(model='llama3.2', messages=[
  {'role': 'user', 'content': 'Hello'}
])

print(response['message']['content'])
```

## Troubleshooting

### "Failed to initialize keystore"

**Linux**: Install libsecret
```bash
sudo apt-get install libsecret-tools
```

**macOS**: Keychain should be available by default. Check macOS version.

**Windows**: Ensure you're running Windows 10 or later.

### "Failed to apply sandbox"

Check platform-specific requirements above and ensure you have necessary permissions.

### "Connection refused" or Network Errors

This is expected! External connections are blocked by design. Only localhost connections are allowed.

### Performance Issues

Encryption adds minimal overhead (~5-10ms per message). If you experience issues:

1. Check system resources (RAM, CPU)
2. Reduce model size
3. Disable encryption temporarily for testing:
```bash
SECLLAMA_ENABLE_ENCRYPTION=false ./secllama serve
```

## Security Maintenance

### Regular Tasks

1. **Update SecLlama regularly**
```bash
git pull origin main
go build .
```

2. **Rotate encryption keys** (monthly recommended)
```bash
./secllama key rotate
```

3. **Review security logs**
```bash
./secllama logs security
```

4. **Audit running processes**
```bash
./secllama ps
# Verify all runners are sandboxed
```

## Uninstallation

### Remove SecLlama

```bash
# Stop all running instances
./secllama stop

# Remove binary
sudo rm /usr/local/bin/secllama

# Remove data (optional)
rm -rf ~/.secllama
rm -rf ~/.ollama  # if you want to remove models too
```

### Remove Encryption Keys

**⚠️ WARNING**: This will make any previously encrypted messages unrecoverable!

```bash
# macOS
security delete-generic-password -s secllama -a message-encryption-key

# Linux
secret-tool clear service secllama account message-encryption-key

# Windows PowerShell
cmdkey /delete:secllama/message-encryption-key
```

### Remove Firewall Rules (Windows)

```bash
# PowerShell as Administrator
netsh advfirewall firewall delete rule name="SecLlama-BlockExternal"
netsh advfirewall firewall delete rule name="SecLlama-BlockExternal-AllowLocal"
```

## Getting Help

- 📖 [Full Security Documentation](SECLLAMA_SECURITY.md)
- 💬 [GitHub Discussions](https://github.com/yourusername/secllama/discussions)
- 🐛 [Report Issues](https://github.com/yourusername/secllama/issues)
- 🔒 [Security Issues](mailto:security@secllama.example.com)

## Next Steps

- Read the [Security Documentation](SECLLAMA_SECURITY.md) to understand the security model
- Check out [CONTRIBUTING.md](CONTRIBUTING.md) to contribute
- Join the community discussions

---

**Note**: SecLlama is designed for maximum security and privacy. Some Ollama features that require external network access (like automatic model downloads during chat, remote models, etc.) are intentionally disabled.

