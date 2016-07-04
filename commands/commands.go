package zpCommands

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type Manager struct {
	db       *bolt.DB
	execList map[string]interface{}
}

type Command struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	Args string `json:"args"`
}

type CommandResult struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type ExecutorInterface interface {
	Init() error
	Exec(command *Command) (string, error)
}

type Executor struct {
	Kind string
}

func NewManager() *Manager {
	newManager := &Manager{}

	newManager.execList = make(map[string]interface{})

	// Create database file
	var err error
	newManager.db, err = bolt.Open("zp.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	return newManager
}

func (m *Manager) Close() {
	m.db.Close()
}

func (m *Manager) Add(command *Command) {

	m.db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte(command.Kind))
		if err != nil {
			return err
		}
		// Encode to json
		encoded, err := json.Marshal(command)
		if err != nil {
			return err
		}
		return b.Put([]byte(command.Name), encoded)
	})

}

func (m *Manager) Get(kind, name string) *Command {

	var command Command

	m.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		v := b.Get([]byte(name))

		if err := json.Unmarshal([]byte(v), &command); err != nil {
			log.Fatal(err)
		}

		return nil
	})
	return &command
}

func (m *Manager) Exec(kind, name, args string) (string, error) {

	// Get the parameters of the command
	command := m.Get(kind, name)
	if len(args) > 0 {
		command.Args = args
	}
	_ = command
	exec := m.execList[kind]
	if exec == nil {

		switch kind {
		case "system":
			execSystem := ExecutorSystem{}
			execSystem.Init()
			m.execList["system"] = execSystem
			return execSystem.Exec(command)
		case "rpi":
			executorRPI := ExecutorRPI{}
			executorRPI.Init()
			m.execList["rpi"] = executorRPI
			return executorRPI.Exec(command)
		case "http":
			executorHTTP := ExecutorHTTP{}
			executorHTTP.Init()
			m.execList["http"] = executorHTTP
			return executorHTTP.Exec(command)
		case "tcp":
			executorTCP := ExecutorTCP{}
			executorTCP.Init()
			m.execList["tcp"] = executorTCP
			return executorTCP.Exec(command)
		}
	} else {
		switch kind {
		case "system":
			exeI := exec.(ExecutorSystem)
			return exeI.Exec(command)
		case "rpi":
			exeRPI := exec.(ExecutorRPI)
			return exeRPI.Exec(command)
		case "http":
			exeHTTP := exec.(ExecutorHTTP)
			return exeHTTP.Exec(command)
		case "tcp":
			exeTCP := exec.(ExecutorTCP)
			return exeTCP.Exec(command)
		}
	}

	return "", nil
}
