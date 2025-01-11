package internal

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"unicode"

	"github.com/fbufler/ssh-tunnel-setup/config"
	"github.com/fbufler/ssh-tunnel-setup/package/sshd"
	"github.com/fbufler/ssh-tunnel-setup/package/system"
)

const minimalPasswordLength = 16

func ServerSetup(cfg *config.ServerConfig) error {
	slog.Info("Setting up server")
	err := setupUser(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to setup user: %v", err))
		return fmt.Errorf("failed to setup user: %v", err)
	}

	err = prepareUserSSHEnvironment(cfg.TunnelUser)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to prepare user's SSH environment: %v", err))
		return fmt.Errorf("failed to prepare user's SSH environment: %v", err)
	}

	err = setupSSHDConfig(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to setup sshd: %v", err))
		return fmt.Errorf("failed to setup sshd: %v", err)
	}

	err = restartSSHD()
	if err != nil {
		slog.Error(fmt.Sprintf("failed to restart sshd: %v", err))
		return fmt.Errorf("failed to restart sshd: %v", err)
	}

	return fmt.Errorf("not implemented")
}

func setupUser(cfg *config.ServerConfig) error {
	slog.Debug("Setting up user")
	if system.UserExists(cfg.TunnelUser) {
		return fmt.Errorf("user already exists, assuming setup was done already")
	}

	if !isPasswordComplex(cfg.TunnelPass) {
		return fmt.Errorf("password does not meet complexity requirements")
	}

	err := system.CreateUser(cfg.TunnelUser, cfg.TunnelPass)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	err = system.AddUserToGroup(cfg.TunnelUser, "tunnel")
	if err != nil {
		return fmt.Errorf("failed to add user to group: %v", err)
	}
	return nil
}

func isPasswordComplex(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= minimalPasswordLength {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func prepareUserSSHEnvironment(user string) error {
	slog.Debug("Preparing user's SSH environment")

	slog.Debug("Creating .ssh directory")
	err := os.Mkdir(fmt.Sprintf("/home/%s/.ssh", user), 0700)
	if err != nil {
		return fmt.Errorf("failed to create .ssh directory: %v", err)
	}
	err = os.Chmod(fmt.Sprintf("/home/%s/.ssh", user), 0700)
	if err != nil {
		return fmt.Errorf("failed to chmod .ssh directory: %v", err)
	}

	slog.Debug("Creating authorized_keys file")
	_, err = os.Create(fmt.Sprintf("/home/%s/.ssh/authorized_keys", user))
	if err != nil {
		return fmt.Errorf("failed to create authorized_keys file: %v", err)
	}
	err = os.Chmod(fmt.Sprintf("/home/%s/.ssh/authorized_keys", user), 0600)
	if err != nil {
		return fmt.Errorf("failed to chmod authorized_keys file: %v", err)
	}

	return nil
}

func setupSSHDConfig(cfg *config.ServerConfig) error {
	slog.Debug("Setting up sshd configuration")
	sshdConfig, err := os.ReadFile(cfg.SSHDConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read sshd config: %v", err)
	}

	sshd.AllowGatewayPorts(&sshdConfig)
	sshd.AllowTcpForwarding(&sshdConfig)
	return os.WriteFile(cfg.SSHDConfigPath, sshdConfig, 0644)
}

func restartSSHD() error {
	slog.Debug("Restarting sshd")
	cmd := exec.Command("sudo", "systemctl", "restart", "sshd")
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to restart sshd: %v", err)
	}
	slog.Debug("sshd restarted")
	return nil
}
