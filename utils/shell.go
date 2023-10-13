package utils

import (
	"bytes"
	"os/exec"
)

func ExecShell(command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
