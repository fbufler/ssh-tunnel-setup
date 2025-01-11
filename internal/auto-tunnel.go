package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/package/ssh"
	"github.com/fbufler/ssh-tunnel-setup/package/system"
)

const serviceName = "managed-tunnel"
const serviceDescription = "Managed SSH tunnel"
const monitorInterval = 60 * time.Second

func SetupTunnel(cfg *config.TunnelConfig) error {
	slog.Debug("Setting up managed tunnel")

	serverKeyPath := cfg.KeyDirectory + "/" + cfg.ServerKeyName
	err := ssh.ConfigureTunnel(cfg.SSHConfigPath, cfg.HostIdentifier, cfg.ServerName, cfg.ServerUser, serverKeyPath, cfg.LocalHost, cfg.LocalPort, cfg.ServerPort)
	if err != nil {
		return err
	}

	execStart := fmt.Sprintf("/usr/bin/ssh -N -R %d:%s:%d %s@%s", cfg.ServerPort, cfg.LocalHost, cfg.LocalPort, cfg.ServerUser, cfg.ServerName)
	err = system.CreateSystemdService(serviceName, serviceDescription, execStart, cfg.LocalUser)
	if err != nil {
		return err
	}

	err = system.EnableSystemdService(serviceName)
	if err != nil {
		return err
	}

	err = system.CreateCronMonitor(serviceName, cfg.LocalUser, cfg.LocalPort, monitorInterval)
	if err != nil {
		return err
	}

	return nil
}
