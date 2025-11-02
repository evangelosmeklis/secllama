<div align="center">
  <img src="secllama_icon.png" alt="SecLlama" width="200">
  <h1>SecLlama 🔒</h1>
</div>

SecLlama is a security-enhanced fork of Ollama that adds critical security layers for running large language models. It provides network isolation through sandboxing (models cannot access the internet), end-to-end encryption for all messages using AES-256-GCM, secure key storage in OS-native keychains, and encrypted command history. All data is stored separately from Ollama in `~/.secllama/` to allow both systems to coexist.

Get started by building with `go build -o secllama .`, pulling a model with `./secllama pull llama3.2`, and running it with `./secllama run llama3.2`. All security features are enabled automatically with no configuration required.
