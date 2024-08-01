package utils

import (
	"fmt"
	"os/exec"
)

type ExecConfig struct {
	Command     string
	WorkingDir  string
	Environment string
}

func ExecCmd(config ExecConfig) (string, error) {
	cmd := exec.Command("sh", "-c", config.Command)

	// configure the env
	cmd.Env = append(cmd.Env, config.Environment)

	// set the working dir (if configured)
	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}

	output, err := cmd.Output()
	if err != nil {
		return string(output), fmt.Errorf("there was an eror executing %s. ERROR: %s", config.Command, err)
	}

	return string(output), nil
}
