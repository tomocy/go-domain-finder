package main

import (
	"log"
	"os"
	"os/exec"
)

var cmdChain = []*exec.Cmd{
	exec.Command("lib/synonyms"),
	exec.Command("lib/sprinkle"),
	exec.Command("lib/coolify"),
	exec.Command("lib/domainify"),
	exec.Command("lib/available"),
}

func main() {
	setUpCommandChain()

	startAndWaitAllCommands()
}

func setUpCommandChain() {
	cmdChain[0].Stdin = os.Stdin
	cmdChain[len(cmdChain)-1].Stdout = os.Stdout
	for i := 0; i < len(cmdChain)-1; i++ {
		cmd := cmdChain[i]
		nextCmd := cmdChain[i+1]
		if err := connectPipe(cmd, nextCmd); err != nil {
			log.Panicf("could not connect pipe: %s\n", err)
		}
	}
}

func connectPipe(cmd *exec.Cmd, nextCmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	nextCmd.Stdin = stdout
	return nil
}

func startAndWaitAllCommands() {
	for _, cmd := range cmdChain {
		if err := cmd.Start(); err != nil {
			log.Panicf("could not start command: %s\n", err)
		}
		defer cmd.Process.Kill()
	}

	for _, cmd := range cmdChain {
		if err := cmd.Wait(); err != nil {
			log.Panicf("could not wait command: %s\n", err)
		}
	}
}
