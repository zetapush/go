package zpCommands

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type ExecutorSystem struct {
	Executor
}

func (e *ExecutorSystem) Init() error {
	return nil
}

func (e *ExecutorSystem) Exec(command *Command) (string, error) {

	log.Printf("ExecutorSystem Exec %#v\n", command)

	var cmd *exec.Cmd

	if len(command.Args) > 0 {
		params := strings.Split(command.Args, ",")
		cmd = exec.Command(command.Cmd, params...)
	} else {
		cmd = exec.Command(command.Cmd)
	}

	// Parse the parameters
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	result := out.String()
	log.Printf("after run %#s\n", result)
	if err != nil {
		log.Printf("Erreur exec %#v\n", err)
		return "", err
	}
	return result, nil
}
