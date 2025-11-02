# SecLlama Quick Start

Get SecLlama running in 5 minutes with maximum security.

## 🚀 Fast Installation

### macOS
```bash
# Clone and build
git clone https://github.com/yourusername/secllama.git
cd secllama
go build .

# Start SecLlama
./secllama serve &

# Pull a model
./secllama pull llama3.2:1b

# Start chatting (encrypted & sandboxed)
./secllama run llama3.2:1b
```

### Linux
```bash
# Install dependencies
sudo apt-get install build-essential util-linux libsecret-tools

# Clone and build
git clone https://github.com/yourusername/secllama.git
cd secllama
go build .

# Start SecLlama
./secllama serve &

# Pull a model
./secllama pull llama3.2:1b

# Start chatting (encrypted & sandboxed)
./secllama run llama3.2:1b
```

### Windows
```bash
# Clone and build (in PowerShell)
git clone https://github.com/yourusername/secllama.git
cd secllama
go build .

# Start SecLlama (requires Administrator for firewall rules)
.\secllama.exe serve

# In another terminal
.\secllama.exe pull llama3.2:1b
.\secllama.exe run llama3.2:1b
```

## 🔒 What's Protected?

Once SecLlama is running:

✅ **All conversations are encrypted** (AES-256-GCM)  
✅ **Models cannot access the internet** (sandboxed)  
✅ **Keys stored securely** (OS keychain)  
✅ **No data leaves your machine** (100% offline)

## 📊 Verify Security

```bash
# Check security status
./secllama security status

# Expected output:
# ✓ Encryption: Enabled (AES-256-GCM)
# ✓ Sandboxing: Enabled
# ✓ Network Isolation: Active
# ✓ Key Storage: Keychain
```

## 🧪 Test Network Isolation

```bash
# Try to make the model access the internet (it will fail)
./secllama run llama3.2:1b "Access httpbin.org/get"

# You should see: Model cannot access external networks
```

## 💡 Usage Examples

### Basic Chat
```bash
./secllama run llama3.2 "What is 2+2?"
```

### Interactive Mode
```bash
./secllama run llama3.2
>>> Hello!
Hello! How can I help you today?
>>> /bye
```

### API Mode
```bash
# Start server
./secllama serve

# Use curl (in another terminal)
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {"role": "user", "content": "Hello!"}
  ]
}'
```

### Python Client
```python
import ollama

client = ollama.Client()
response = client.chat(model='llama3.2', messages=[
  {'role': 'user', 'content': 'Hello!'}
])
print(response['message']['content'])
```

## 🔧 Configuration

SecLlama works out of the box with secure defaults. To customize:

```bash
# Disable encryption (not recommended)
export SECLLAMA_ENABLE_ENCRYPTION=false

# Disable sandboxing (not recommended)
export SECLLAMA_ENABLE_SANDBOX=false

# Check configuration
./secllama config show
```

## 📚 Learn More

- **Full Setup Guide**: [SETUP.md](SETUP.md)
- **Security Details**: [SECLLAMA_SECURITY.md](SECLLAMA_SECURITY.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)

## ❓ Common Issues

### "Failed to initialize keystore"
Install platform-specific tools:
- **Linux**: `sudo apt-get install libsecret-tools`
- **macOS**: Should work out of the box
- **Windows**: Should work out of the box

### "Failed to apply sandbox"
- **Linux**: Install `util-linux` package
- **macOS**: Should work out of the box
- **Windows**: Run as Administrator

### Models are slow
SecLlama has minimal overhead (~5-10ms). If slow:
1. Use smaller models (try `llama3.2:1b`)
2. Check system resources
3. Ensure GPU drivers are installed

## 🆘 Get Help

- 📖 [Documentation](SECLLAMA_SECURITY.md)
- 💬 [GitHub Discussions](https://github.com/yourusername/secllama/discussions)
- 🐛 [Report Issues](https://github.com/yourusername/secllama/issues)

---

**Ready for more?** Check out the [full setup guide](SETUP.md) for advanced configuration options!

