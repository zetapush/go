package zpCommands

import (
	"io/ioutil"
	"log"
	"net/http"
)

type ExecutorHTTP struct {
	Executor
}

func (e *ExecutorHTTP) Init() error {

	log.Println("ExecutorHTTP Init")
	// nothing here
	return nil
}

func (e *ExecutorHTTP) Exec(command *Command) (string, error) {

	log.Println("ExecutorHTTP Exec")

	// TODO Better error handling

	switch command.Cmd {
	case "GET":
		res, err := http.Get(command.Args)
		if err != nil {
			return "", err
		}
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return "", err
		}
		return string(result), err
	}

	return "", nil

}
