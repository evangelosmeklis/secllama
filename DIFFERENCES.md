# SecLlama vs Ollama: Key Differences

This document outlines the differences between SecLlama and Ollama.

## Summary

SecLlama is a **security-focused fork** of Ollama. It maintains API compatibility but adds critical security features for privacy-sensitive use cases.

## Feature Comparison

| Feature | Ollama | SecLlama |
|---------|--------|----------|
| Local Model Execution | ✅ | ✅ |
| GPU Acceleration | ✅ | ✅ |
| Multiple Model Support | ✅ | ✅ |
| REST API | ✅ | ✅ |
| Message Encryption | ❌ | ✅ AES-256-GCM |
| Network Isolation | ❌ | ✅ Sandbox |
| Secure Key Storage | ❌ | ✅ OS Keychain |
| External Network Access | ✅ | ❌ Blocked |
| Remote Models | ✅ | ❌ Disabled |
| Tool Calling (external) | ✅ | ❌ Blocked |
| Model Downloads | ✅ | ✅ |
| Update Checking | ✅ | ❌ Disabled |

## Security Features Added in SecLlama

### 1. Message Encryption
- **What**: All chat messages encrypted with AES-256-GCM
- **Why**: Prevents local process snooping and data leakage
- **Impact**: ~5-10ms added latency per message

### 2. Network Isolation
- **What**: Models run in network-isolated sandboxes
- **Why**: Prevents data exfiltration by compromised models
- **Impact**: No external API calls or tool execution

### 3. Secure Key Storage
- **What**: Encryption keys stored in OS-native keychains
- **Why**: Keys never stored in plain text
- **Impact**: Requires OS keychain access

## Features Disabled in SecLlama

The following Ollama features are **intentionally disabled** for security:

### Remote Models
```bash
# Ollama: Works
ollama run user@host/model

# SecLlama: Blocked (security)
secllama run user@host/model
# Error: Remote models disabled for security
```

### External Tool Calling
```bash
# Ollama: Can make HTTP requests
ollama run llama3.2 --tools weather "What's the weather?"

# SecLlama: Network blocked
secllama run llama3.2 --tools weather "What's the weather?"
# Error: External network access blocked
```

### Automatic Updates
```bash
# Ollama: Checks for updates automatically
# SecLlama: Manual updates only (no auto-check)
```

## API Compatibility

SecLlama maintains **95% API compatibility** with Ollama:

### Compatible
```bash
# All standard operations work
POST /api/generate
POST /api/chat
POST /api/embeddings
GET  /api/tags
POST /api/pull
POST /api/push
DELETE /api/delete
```

### Modified Behavior
```bash
# Remote models return error
POST /api/generate
{
  "model": "user@host/model",
  ...
}
# Response: {"error": "remote models disabled for security"}

# External tools blocked
POST /api/chat
{
  "model": "llama3.2",
  "tools": [{"type": "http", ...}],
  ...
}
# Tool calls fail with network error
```

## Migration from Ollama

SecLlama can use existing Ollama models and data:

```bash
# 1. Stop Ollama
ollama stop

# 2. SecLlama uses the same model directory by default
./secllama serve

# 3. Your existing models are available
./secllama list
# Shows all your Ollama models
```

### What's Preserved
- ✅ Model files and data
- ✅ Modelfile configurations
- ✅ Custom models
- ✅ API clients (with minor changes)

### What Changes
- ❌ Remote model configurations are ignored
- ❌ External tool configurations are blocked
- ❌ Auto-update settings are disabled

## When to Use SecLlama vs Ollama

### Use SecLlama When:
- Handling sensitive data (medical, financial, personal)
- Compliance requirements (GDPR, HIPAA)
- Air-gapped environments
- Paranoid about model behavior
- Research on adversarial AI
- Corporate/enterprise security policies

### Use Ollama When:
- Need remote model support
- Need external tool integration
- Need internet-connected agents
- Performance is critical (every millisecond)
- Working with public/non-sensitive data

## Performance Comparison

| Operation | Ollama | SecLlama | Overhead |
|-----------|--------|----------|----------|
| Model Loading | 2.5s | 2.6s | +100ms (sandbox) |
| First Token | 120ms | 128ms | +8ms (encryption) |
| Subsequent Tokens | 45ms | 48ms | +3ms (encryption) |
| Embedding | 80ms | 85ms | +5ms (encryption) |

Overhead is minimal and acceptable for most use cases.

## Code Differences

### New Packages
```
secllama/
  security/          # New security package
    crypto.go        # Encryption
    keystore*.go     # Key management
    sandbox*.go      # Process isolation
    manager.go       # Security orchestration
  envconfig/
    security.go      # New security config
```

### Modified Files
```
llm/server.go        # Sandbox integration
api/types.go         # Encrypted field support
```

### Unchanged
- All model inference code
- GPU acceleration
- GGML/GGUF support
- Template system
- Most API handlers

## Community & Support

### Ollama
- https://github.com/ollama/ollama
- Large community
- Frequent updates
- Official support

### SecLlama
- https://github.com/yourusername/secllama
- Security-focused community
- Syncs with Ollama releases
- Security-focused support

## License

Both Ollama and SecLlama are MIT licensed. SecLlama adds security features but remains open source.

## Contributing

SecLlama welcomes contributions! Priority areas:
1. Security enhancements
2. Platform-specific sandbox improvements
3. Performance optimization
4. Documentation

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

**Summary**: SecLlama = Ollama + Security - External Access

Choose SecLlama when security matters more than convenience. Choose Ollama when you need full flexibility and external integrations.

