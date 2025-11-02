package envconfig

import (
	"os"
	"strconv"
)

// EnableEncryption returns whether message encryption is enabled
func EnableEncryption() bool {
	if enabled := os.Getenv("SECLLAMA_ENABLE_ENCRYPTION"); enabled != "" {
		val, _ := strconv.ParseBool(enabled)
		return val
	}
	// Default to enabled for secllama
	return true
}

// EnableSandbox returns whether sandboxing is enabled
func EnableSandbox() bool {
	if enabled := os.Getenv("SECLLAMA_ENABLE_SANDBOX"); enabled != "" {
		val, _ := strconv.ParseBool(enabled)
		return val
	}
	// Default to enabled for secllama
	return true
}

// StrictNetworkIsolation returns whether strict network isolation is enforced
func StrictNetworkIsolation() bool {
	if enabled := os.Getenv("SECLLAMA_STRICT_NETWORK_ISOLATION"); enabled != "" {
		val, _ := strconv.ParseBool(enabled)
		return val
	}
	// Default to enabled for secllama
	return true
}

