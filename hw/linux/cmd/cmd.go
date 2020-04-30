package cmd

import (
	"os/exec"

	"log"
)

// Exec Execute a command and collect the output
func Exec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Printf("Exec: %s %s", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)
	res, err := cmd.CombinedOutput()

	return string(res), err
}
