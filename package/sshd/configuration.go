package sshd

import (
	"log/slog"
	"regexp"
)

// AllowGatewayPorts sets GatewayPorts to yes in the sshd_config file
func AllowGatewayPorts(sshdConfig *[]byte) {

	gatewayPortsPattern := regexp.MustCompile(`(?m)^\s*#?\s*GatewayPorts\s+yes\s*$`)
	if !gatewayPortsPattern.Match(*sshdConfig) {
		slog.Debug("Adding GatewayPorts yes to sshd_config")
		*sshdConfig = append(*sshdConfig, []byte("\nGatewayPorts yes\n")...)
	} else {
		slog.Debug("Uncommenting GatewayPorts yes in sshd_config")
		*sshdConfig = gatewayPortsPattern.ReplaceAll(*sshdConfig, []byte("GatewayPorts yes"))
	}
}

// AllowTcpForwarding sets AllowTcpForwarding to yes in the sshd_config file
func AllowTcpForwarding(sshdConfig *[]byte) {
	allowTcpForwardingPattern := regexp.MustCompile(`(?m)^\s*#?\s*AllowTcpForwarding\s+yes\s*$`)
	if !allowTcpForwardingPattern.Match(*sshdConfig) {
		slog.Debug("Adding AllowTcpForwarding yes to sshd_config")
		*sshdConfig = append(*sshdConfig, []byte("\nAllowTcpForwarding yes\n")...)
	} else {
		slog.Debug("Uncommenting AllowTcpForwarding yes in sshd_config")
		*sshdConfig = allowTcpForwardingPattern.ReplaceAll(*sshdConfig, []byte("AllowTcpForwarding yes"))
	}
}
