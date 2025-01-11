package internal

import (
	"fmt"
	"log/slog"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/package/ssh"
)

func Rotate(cfg *config.RotateConfig) error {
	slog.Info("Rotating key pair")
	serverAdress := fmt.Sprintf("%s:%d", cfg.ServerName, cfg.ServerPort)
	err := ssh.RotateKeyPair(cfg.ServerUser, serverAdress, cfg.KeyDirectory, cfg.KeyName, cfg.KeyUser, config.TrustedHostKey())
	if err != nil {
		slog.Error(fmt.Sprintf("Error rotating key pair: %s", err))
		return err
	}
	slog.Info("Key pair rotated")

	slog.Info("Update Rotation Config")
	err = config.StoreRotationConfig(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error adding rotation config: %s", err))
		return err
	}

	return nil
}
