package system

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

const systemdDir = "/etc/systemd/system/"

func CreateSystemdService(serviceName, description, execStart, user string) error {
	slog.Debug("Creating systemd service")

	slog.Debug("Verifying systemd directory")
	_, err := os.Stat(systemdDir)
	if err != nil {
		return fmt.Errorf("failed to check systemd directory: %v", err)
	}

	servicePath := systemdDir + serviceName + ".service"
	slog.Debug("Creating service file")
	serviceFile, err := os.Create(servicePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %v", err)
	}

	defer serviceFile.Close()

	slog.Debug("Writing service configuration")
	serviceConfig := fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
ExecStart=%s
Restart=always
User=%s

[Install]
WantedBy=multi-user.target
`, description, execStart, user)

	_, err = serviceFile.WriteString(serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to write service configuration: %v", err)
	}

	slog.Debug("Setting service file permissions")
	err = os.Chmod(servicePath, 0644)
	if err != nil {
		return fmt.Errorf("failed to chmod service file: %v", err)
	}

	slog.Debug("Setting service user")
	err = exec.Command("chown", user, servicePath).Run()
	if err != nil {
		return fmt.Errorf("failed to set service user: %v", err)
	}

	return nil
}

func EnableSystemdService(serviceName string) error {
	slog.Debug("Enabling systemd service")

	slog.Debug("Reloading systemd")
	err := exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}

	slog.Debug("Enabling service")
	err = exec.Command("systemctl", "enable", serviceName).Run()
	if err != nil {
		return fmt.Errorf("failed to enable service: %v", err)
	}

	return nil
}
