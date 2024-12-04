package main

import (
	// "io".
	"log"
	"os"
	"os/exec"
)

func prepareEnv(env Environment) error {
	for k, v := range env {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				return err
			}
			continue
		}
		if err := os.Setenv(k, v.Value); err != nil {
			return err
		}
	}
	return nil
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	prepareEnv(env)
	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := command.Start(); err != nil {
		log.Fatalf("unable to start program, %v", err)
	}
	if err := command.Wait(); err != nil {
		switch e := err.(type) { //nolint:errorlint
		case *exec.ExitError:
			return e.ExitCode()
		default:
			log.Fatalf("unexpected error during program execution, %v", e)
			return
		}
	}
	return 0
}
