// utils/shell.go

package utils

import (
	"bytes"
	"os/exec"
)

// ExecuteShellCommand executes a shell command and returns its output as a string.
func ExecuteShellCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
