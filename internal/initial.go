package internal

import (
	"fmt"
	"log/slog"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/package/ssh"
)

func ClientSetup(cfg *config.ClientConfig) error {
	slog.Info("Setting up client")

	err := ssh.PrepareKeyDirectory(cfg.KeyDirectory)
	if err != nil {
		slog.Error(fmt.Sprintf("Error preparing key directory: %s", err))
		return err
	}

	slog.Info("Generating key pair")
	err = ssh.MakeKeyPair(cfg.KeyDirectory, cfg.KeyName, cfg.KeyUser)
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating key pair: %s", err))
		return err
	}
	slog.Info(fmt.Sprintf("Key pair generated at %s/%s", cfg.KeyDirectory, cfg.KeyName))

	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerName, cfg.ServerPort)

	slog.Info(fmt.Sprintf("Authorizing public key on remote: %s", serverAddr))
	remoteAuth := ssh.RemoteAuth{
		User: cfg.ServerUser,
	}
	if cfg.ServerPass != "" {
		remoteAuth.Password = cfg.ServerPass
	}

	if cfg.ServerKeyName != "" {
		remoteAuth.KeyPath = fmt.Sprintf("%s/%s", cfg.KeyDirectory, cfg.ServerKeyName)
	}

	err = ssh.AuthorizePublicKeyOnRemote(fmt.Sprintf("%s/%s", cfg.KeyDirectory, cfg.KeyName), serverAddr, config.TunnelUser, remoteAuth)
	if err != nil {
		slog.Error(fmt.Sprintf("Error authorizing public key on remote: %s", err))
		return err
	}
	slog.Info("Public key authorized on remote")

	// TODO: Add automatic tunnel setup

	slog.Info("Add Rotation Config")
	err = addRotationConfig(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error adding rotation config: %s", err))
		return err
	}

	return nil
}

func addRotationConfig(cfg *config.ClientConfig) error {
	newRotationConfig := config.RotateConfig{
		KeyName:      cfg.KeyName,
		KeyDirectory: cfg.KeyDirectory,
		KeyUser:      cfg.KeyUser,
		ServerName:   cfg.ServerName,
		ServerPort:   cfg.ServerPort,
		ServerUser:   config.TunnelUser,
	}

	return config.StoreRotationConfig(&newRotationConfig)
}
