package zpCommands

import (
	"log"
)

type ExecutorRPI struct {
	Executor
}

func (e *ExecutorRPI) Init() error {

	log.Println("ExecutorRPI Init")

	return nil
}

func (e *ExecutorRPI) Exec(command *Command) (string, error) {

	log.Println("ExecutorRPI Exec")
	return "", nil
}
