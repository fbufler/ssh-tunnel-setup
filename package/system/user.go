package system

import (
	"fmt"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func checkOS(assumed string) bool {
	currentOS := runtime.GOOS
	return currentOS == assumed
}

func checkSudo() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}

	currentUserGroups, err := currentUser.GroupIds()
	if err != nil {
		return false
	}
	for _, groupID := range currentUserGroups {
		systemGroup, err := user.LookupGroupId(groupID)
		if err != nil {
			continue
		}
		if systemGroup.Name == "sudo" {
			return true
		}
	}
	return false
}

func UserExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

// CreateUser creates a new user with the provided username and password
func CreateUser(username, password string) error {
	if !checkOS("linux") {
		return fmt.Errorf("only supported on Linux")
	}

	if !checkSudo() {
		return fmt.Errorf("root access is required for user creation")
	}

	if UserExists(username) {
		return fmt.Errorf("user %s already exists", username)
	}

	cmd := exec.Command("sudo", "useradd", "-m", "-p", password, username)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	cmd = exec.Command("sudo", "chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", username, password))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to set user password: %v", err)
	}

	return nil
}

func groupExists(group string) bool {
	_, err := user.LookupGroup(group)
	return err == nil
}

func AddUserToGroup(username, group string) error {
	if !checkOS("linux") {
		return fmt.Errorf("only supported on Linux")
	}

	if !checkSudo() {
		return fmt.Errorf("root access is required for user management")
	}

	if !UserExists(username) {
		return fmt.Errorf("user %s does not exist", username)
	}

	if !groupExists(group) {
		return fmt.Errorf("group %s does not exist", group)
	}

	cmd := exec.Command("sudo", "usermod", "-aG", group, username)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add user to group: %v", err)
	}

	return nil
}

func HomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
