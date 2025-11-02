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
	
	// Allow most operations by default, only restrict specific things
	sb.WriteString("(allow default)\n\n")
	
	// Deny external network access if needed
	// Note: macOS sandbox-exec has limited fine-grained network filtering
	// The main network security is enforced at the application layer via
	// the localhost-only HTTP client. The sandbox provides process isolation.
	if !config.AllowLocalhost {
		sb.WriteString("; Deny all network access\n")
		sb.WriteString("(deny network-outbound)\n")
		sb.WriteString("(deny network-bind)\n\n")
	}
	
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

