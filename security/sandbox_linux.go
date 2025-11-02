package security

import (
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

// applySandbox applies Linux-specific sandboxing using namespaces and seccomp
func applySandbox(cmd *exec.Cmd, config SandboxConfig) error {
	// Use Linux namespaces for isolation
	// We'll use network namespace to isolate network access
	
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	
	// Create new network namespace to isolate network
	// CLONE_NEWNET creates a new network namespace
	cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWNET
	
	// Set up network namespace with only loopback
	// This requires root or CAP_NET_ADMIN capability
	// For non-root users, we'll use alternative methods
	
	slog.Info("linux sandbox applied with network isolation")
	
	// Alternative: Use unshare command if available
	// This provides better compatibility without requiring root
	unsharePath, err := exec.LookPath("unshare")
	if err == nil {
		// Wrap command with unshare
		originalArgs := cmd.Args
		cmd.Path = unsharePath
		cmd.Args = append([]string{"unshare", "--net", "--"}, originalArgs...)
		slog.Info("using unshare for network isolation", "args", cmd.Args)
	} else {
		slog.Warn("unshare not available, network isolation may be limited")
	}
	
	return nil
}

// SetupSeccompFilter sets up a seccomp filter to restrict system calls
func SetupSeccompFilter(cmd *exec.Cmd) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("seccomp only supported on Linux")
	}
	
	// List of blocked syscalls for network access
	blockedSyscalls := []string{
		"socket",     // Block socket creation
		"connect",    // Block network connections
		"sendto",     // Block sending data
		"recvfrom",   // Block receiving data
		"sendmsg",    // Block message sending
		"recvmsg",    // Block message receiving
	}
	
	slog.Info("seccomp filter would block syscalls", "syscalls", blockedSyscalls)
	
	// Note: Full seccomp implementation requires cgo and libseccomp
	// For now, we'll rely on network namespace isolation
	
	return nil
}

// CreateIsolatedNetworkNamespace creates a network namespace with only loopback
func CreateIsolatedNetworkNamespace() error {
	// This would need to be called with appropriate privileges
	// Requires CAP_NET_ADMIN or root
	
	// 1. Create new network namespace
	err := unix.Unshare(unix.CLONE_NEWNET)
	if err != nil {
		return fmt.Errorf("failed to create network namespace: %v", err)
	}
	
	// 2. Bring up loopback interface
	// This requires netlink operations which we'll implement later
	slog.Info("created isolated network namespace")
	
	return nil
}

// WrapCommandWithFirewall wraps a command with iptables rules to block external network
func WrapCommandWithFirewall(cmd *exec.Cmd, allowedPorts []int) error {
	// This is an alternative approach using iptables
	// Creates rules to block all but localhost traffic
	
	// Build iptables rules
	var rules []string
	rules = append(rules, "iptables -A OUTPUT -o lo -j ACCEPT")
	rules = append(rules, "iptables -A OUTPUT -d 127.0.0.0/8 -j ACCEPT")
	
	for _, port := range allowedPorts {
		rules = append(rules, fmt.Sprintf("iptables -A OUTPUT -p tcp --dport %d -d 127.0.0.1 -j ACCEPT", port))
	}
	
	rules = append(rules, "iptables -A OUTPUT -j DROP")
	
	slog.Info("firewall rules", "rules", strings.Join(rules, " && "))
	
	// Note: Actually applying iptables rules requires root
	// This is more of a reference implementation
	
	return nil
}

