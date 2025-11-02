package security

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// applySandbox applies macOS-specific sandboxing using sandbox-exec
func applySandbox(cmd *exec.Cmd, config SandboxConfig) error {
	// macOS has built-in sandboxing via sandbox-exec
	// We'll create a sandbox profile that blocks network access
	
	profile := generateSandboxProfile(config)
	
	// Write profile to temporary file
	tmpDir := os.TempDir()
	profilePath := filepath.Join(tmpDir, "secllama-sandbox.sb")
	
	err := os.WriteFile(profilePath, []byte(profile), 0600)
	if err != nil {
		return fmt.Errorf("failed to write sandbox profile: %v", err)
	}
	
	// Wrap command with sandbox-exec
	originalPath := cmd.Path
	originalArgs := cmd.Args
	
	cmd.Path = "/usr/bin/sandbox-exec"
	cmd.Args = append([]string{"sandbox-exec", "-f", profilePath, originalPath}, originalArgs[1:]...)
	
	slog.Info("macOS sandbox applied", "profile", profilePath)
	
	return nil
}

// generateSandboxProfile creates a macOS sandbox profile that blocks network
func generateSandboxProfile(config SandboxConfig) string {
	var sb strings.Builder
	
	sb.WriteString("(version 1)\n")
	sb.WriteString("(debug deny)\n\n")
	
	// Deny all by default
	sb.WriteString("(deny default)\n\n")
	
	// Allow basic operations
	sb.WriteString("; Allow reading\n")
	sb.WriteString("(allow file-read*)\n\n")
	
	sb.WriteString("; Allow writing to temp and specified paths\n")
	sb.WriteString("(allow file-write*\n")
	sb.WriteString("  (subpath \"/tmp\")\n")
	sb.WriteString("  (subpath \"/private/tmp\")\n")
	sb.WriteString("  (subpath \"/var/tmp\")\n")
	
	for _, path := range config.AllowedWritePaths {
		sb.WriteString(fmt.Sprintf("  (subpath \"%s\")\n", path))
	}
	sb.WriteString(")\n\n")
	
	// Allow process operations
	sb.WriteString("; Allow process operations\n")
	sb.WriteString("(allow process*)\n")
	sb.WriteString("(allow signal)\n")
	sb.WriteString("(allow sysctl-read)\n")
	sb.WriteString("(allow mach-lookup)\n\n")
	
	// Block all network except localhost if configured
	if config.AllowLocalhost {
		sb.WriteString("; Allow localhost connections only\n")
		sb.WriteString("(allow network*\n")
		sb.WriteString("  (remote ip \"localhost:*\")\n")
		sb.WriteString("  (remote ip \"127.0.0.1:*\")\n")
		
		for _, port := range config.AllowedPorts {
			sb.WriteString(fmt.Sprintf("  (remote ip \"127.0.0.1:%d\")\n", port))
		}
		
		sb.WriteString(")\n\n")
	} else {
		sb.WriteString("; Deny all network access\n")
		sb.WriteString("(deny network*)\n\n")
	}
	
	// Explicitly deny internet access
	sb.WriteString("; Explicitly deny external network\n")
	sb.WriteString("(deny network-outbound)\n")
	sb.WriteString("(deny network-bind)\n")
	sb.WriteString("(deny network-inbound)\n")
	
	return sb.String()
}

// CheckSandboxSupport checks if sandbox-exec is available on macOS
func CheckSandboxSupport() error {
	_, err := exec.LookPath("sandbox-exec")
	if err != nil {
		return fmt.Errorf("sandbox-exec not found on macOS")
	}
	return nil
}

