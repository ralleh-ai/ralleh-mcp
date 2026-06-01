package netguard

import (
	"fmt"
	"net"
	"strings"
)

// ValidateLocalListen rejects non-loopback listen addresses unless explicitly
// allowed. Ralleh MCP is intended to be private to the VPS/OpenClaw runtime.
func ValidateLocalListen(addr string, allowNonLoopback bool) error {
	if addr == "" {
		return fmt.Errorf("listen address is required")
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid listen address %q: %w", addr, err)
	}
	if host == "" {
		host = "0.0.0.0"
	}
	if allowNonLoopback {
		return nil
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("listen host %q must be an IP or localhost", host)
	}
	if !ip.IsLoopback() {
		return fmt.Errorf("refusing non-loopback listen address %q; use explicit override only with firewall/auth controls", addr)
	}
	return nil
}
