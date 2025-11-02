# SecLlama System Tray App

This is a native system tray application for running SecLlama with full security features.

## Features

- 🔒 **Secure by Default** - All security features enabled automatically
- 🎯 **Lightweight** - Native Go implementation
- 🔔 **System Tray** - Quick access from your taskbar/menu bar
- ⚡ **Auto-start** - Optional boot-time startup
- 📊 **Status Monitoring** - Real-time server status

This provides the same secure, encrypted, sandboxed AI experience as the CLI, with a convenient GUI.

## Building

### Linux

```bash
cd app
go build -o secllama-app .
./secllama-app
```

### macOS

```bash
cd app
go build -o secllama-app .
./secllama-app
```

### Windows

If you want to build the installer, you'll need to install:
- https://jrsoftware.org/isinfo.php

In the top directory of this repo, run the following powershell script
to build the SecLlama CLI, tray app, and installer:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\build_windows.ps1
```
