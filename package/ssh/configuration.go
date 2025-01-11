package ssh

import (
	"log/slog"
	"os"
	"regexp"
)

// ConfigureTunnel configures an SSH tunnel in the ssh config file
func ConfigureTunnel(sshConfigPath, hostIdentifier, hostName, user, identityFile, localHost string, localPort, remotePort int) error {
	slog.Debug("Configuring tunnel")
	sshConfig, err := os.ReadFile(sshConfigPath)
	if err != nil {
		return err
	}

	slog.Debug("Checking if host is already configured")
	if hostConfigured(sshConfig, hostIdentifier) {
		return nil
	}

	slog.Debug("Adding host configuration")
	hostConfig := []byte("\nHost " + hostIdentifier + "\n")
	hostConfig = append(hostConfig, []byte("    HostName "+hostName+"\n")...)
	hostConfig = append(hostConfig, []byte("    User "+user+"\n")...)
	hostConfig = append(hostConfig, []byte("    IdentityFile "+identityFile+"\n")...)
	hostConfig = append(hostConfig, []byte("    RemoteForward "+localHost+":"+string(localPort)+" localhost:"+string(remotePort)+"\n")...)
	sshConfig = append(sshConfig, hostConfig...)

	slog.Debug("Writing ssh config")
	err = os.WriteFile(sshConfigPath, sshConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func hostConfigured(sshConfig []byte, hostIdentifier string) bool {
	hostPattern := `(?m)^\s*Host\s+` + hostIdentifier + `\s*$`
	return regexp.MustCompile(hostPattern).Match(sshConfig)
}
