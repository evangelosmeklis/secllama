package security

import (
	"log/slog"
	"os/exec"
)

// SandboxConfig defines configuration for sandboxing runner processes
type SandboxConfig struct {
	// AllowLocalhost allows connections to localhost (required for runner communication)
	AllowLocalhost bool
	// AllowedPorts specifies which ports can be accessed (if AllowLocalhost is true)
	AllowedPorts []int
	// WorkingDirectory is the directory where the runner can operate
	WorkingDirectory string
	// AllowedReadPaths are paths the runner can read from
	AllowedReadPaths []string
	// AllowedWritePaths are paths the runner can write to
	AllowedWritePaths []string
}

// ApplySandbox applies OS-specific sandboxing to a command
// This is implemented in platform-specific files (sandbox_*.go)
func ApplySandbox(cmd *exec.Cmd, config SandboxConfig) error {
	slog.Info("applying sandbox to runner process", "config", config)
	return applySandbox(cmd, config)
}

// NetworkIsolationLevel defines how strictly network access is restricted
type NetworkIsolationLevel int

const (
	// NoNetwork completely blocks all network access
	NoNetwork NetworkIsolationLevel = iota
	// LocalhostOnly allows only localhost connections
	LocalhostOnly
	// LocalNetworkOnly allows local network (LAN) connections
	LocalNetworkOnly
)

