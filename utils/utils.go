package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type ExecConfig struct {
	Command     string
	Args        []string
	WorkingDir  string
	Environment []string
}

func ExecCmd(config ExecConfig) (string, error) {
	cmd := exec.Command(config.Command, config.Args...)

	output := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = output

	// configure the env
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, config.Environment...)

	// set the working dir (if configured)
	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}

	err := cmd.Run()
	if err != nil {
		return output.String(), fmt.Errorf("there was an eror executing %s. ERROR: %s", config.Command, err)
	}

	return output.String(), nil
}
