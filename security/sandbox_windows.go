package security

import (
	"fmt"
	"log/slog"
	"os/exec"
)

// applySandbox applies Windows-specific sandboxing
func applySandbox(cmd *exec.Cmd, config SandboxConfig) error {
	// Windows sandboxing options:
	// 1. Use Windows Firewall rules
	// 2. Use AppContainer (requires more complex setup)
	// 3. Use Job Objects with limited permissions
	
	slog.Warn("Windows sandboxing has limited support - applying firewall rules")
	
	// For Windows, we'll create firewall rules to block the process
	// This requires admin privileges
	
	// Alternative: Set restricted token for the process
	// This is more complex and requires Win32 API calls
	
	err := createWindowsFirewallRule(config)
	if err != nil {
		slog.Warn("failed to create firewall rule", "error", err)
		// Don't fail - continue with limited isolation
	}
	
	return nil
}

// createWindowsFirewallRule creates a Windows Firewall rule to block external network
func createWindowsFirewallRule(config SandboxConfig) error {
	ruleName := "SecLlama-BlockExternal"
	
	// Remove existing rule if it exists
	removeCmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		fmt.Sprintf("name=%s", ruleName))
	_ = removeCmd.Run()
	
	// Create new rule to block all outbound traffic
	// Exception: Allow localhost traffic
	addCmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s", ruleName),
		"dir=out",
		"action=block",
		"enable=yes",
		"profile=any")
	
	output, err := addCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create firewall rule: %v, output: %s", err, string(output))
	}
	
	slog.Info("created Windows firewall rule", "rule", ruleName)
	
	// Add exception for localhost
	allowLocalCmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		fmt.Sprintf("name=%s-AllowLocal", ruleName),
		"dir=out",
		"action=allow",
		"enable=yes",
		"remoteip=127.0.0.1")
	
	_, _ = allowLocalCmd.CombinedOutput()
	
	return nil
}

// RemoveWindowsFirewallRule removes the SecLlama firewall rule
func RemoveWindowsFirewallRule() error {
	ruleName := "SecLlama-BlockExternal"
	
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		fmt.Sprintf("name=%s", ruleName))
	
	_ = cmd.Run()
	
	// Also remove localhost exception
	cmd2 := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		fmt.Sprintf("name=%s-AllowLocal", ruleName))
	
	_ = cmd2.Run()
	
	return nil
}

// UseJobObjects configures Windows Job Objects for process isolation
func UseJobObjects(cmd *exec.Cmd) error {
	// This would require Win32 API calls via cgo or syscall
	// Job Objects can limit network access, CPU, memory, etc.
	// This is a placeholder for future implementation
	
	slog.Info("Windows Job Objects not yet implemented")
	return nil
}

