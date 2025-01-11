package system

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

/*
### Monitor Tunnel Uptime with Cron

1. **Create a Monitoring Script:**

		```bash
		nano /path/to/tunnel-check.sh
		```

		Add the following:

		```bash
		#!/bin/bash
		if ! nc -z localhost 2222; then
		    systemctl restart ssh-tunnel
		fi
		```

	 2. **Schedule the Script Using Cron:**
	    ```bash
	    crontab -e
	    ```
	    Add the following line:
	    ```text
	    /5 *
*/
const crontabDir = "/var/spool/cron/crontabs"

func CreateCronMonitor(serviceName, user string, localPort int, interval time.Duration) error {
	slog.Debug("Creating cron monitor")

	err := ensureCrontabDir(user)
	if err != nil {
		return fmt.Errorf("failed to ensure crontab directory: %v", err)
	}

	cronPath := fmt.Sprintf("%s/%s", crontabDir, user)
	slog.Debug("Creating cron file")
	cronFile, err := os.Create(cronPath)
	if err != nil {
		return fmt.Errorf("failed to create cron file: %v", err)
	}

	defer cronFile.Close()

	scriptPath := createServiceMonitor(serviceName, user, localPort, "/usr/local/bin")
	if scriptPath == "" {
		return fmt.Errorf("failed to create monitor script")
	}

	slog.Debug("Writing cron configuration")
	cronConfig := fmt.Sprintf("*/%d * * * * %s\n", interval, scriptPath)
	_, err = cronFile.WriteString(cronConfig)
	if err != nil {
		return fmt.Errorf("failed to write cron configuration: %v", err)
	}

	slog.Debug("Setting cron file permissions")
	err = os.Chmod(cronPath, 0600)
	if err != nil {
		return fmt.Errorf("failed to chmod cron file: %v", err)
	}

	return nil
}

func ensureCrontabDir(user string) error {
	slog.Debug("Ensuring crontab directory in user space")
	err := os.MkdirAll("/var/spool/cron/crontabs", 0700)
	if err != nil {
		return fmt.Errorf("failed to create crontab directory: %v", err)
	}

	return nil
}

func createServiceMonitor(serviceName, user string, localPort int, scriptDir string) string {
	scriptPath := fmt.Sprintf("%s/%s-monitor.sh", scriptDir, serviceName)
	script := fmt.Sprintf(`#!/bin/bash
if ! nc -z localhost %d; then
	systemctl restart %s
fi
`, localPort, serviceName)

	err := os.WriteFile(scriptPath, []byte(script), 0700)
	if err != nil {
		return ""
	}

	return scriptPath
}
