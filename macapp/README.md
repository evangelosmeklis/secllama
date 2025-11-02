# SecLlama Desktop App

This app provides a desktop experience for running SecLlama with full security features.

## Features

- 🔒 **Encrypted Communication** - All messages encrypted with AES-256-GCM
- 🏝️ **Network Isolated** - Models run in sandboxes with no internet access
- 🔑 **Secure Key Storage** - Keys stored in OS-native keychains
- 💻 **Native GUI** - Beautiful desktop interface for macOS/Windows/Linux

## Developing

First, build the `secllama` binary:

```shell
cd ..
go build .
```

Then run the desktop app with `npm start`:

```shell
cd macapp
npm install
npm start
```

The app will automatically start the SecLlama server in the background with all security features enabled.

