package builder

import (
	"os"
	"os/exec"
)

func newCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if os.Getenv("DEBUG") != "false" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
