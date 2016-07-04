package zpCommands

import (
	"log"
	"fmt"
	"net"
	"bufio"
)
var conn net.Conn

type ExecutorTCP struct {
	Executor
}

func (e *ExecutorTCP) Init() error {

	log.Println("ExecutorHTTP Init")
	// nothing here
	return nil
}

func (e *ExecutorTCP) Exec(command *Command) (string, error) {

	log.Println("ExecutorTCP Exec")

	switch command.Cmd {
	case "CONNECT":
		var err error
		conn, err = net.Dial("tcp", command.Args)
				
		if err != nil {
			return "", err
		}
		return "ok", nil
	case "DISCONNECT":
		conn.Close()
		conn= nil	
		return "ok", nil
	case "SEND":
		if conn!= nil{
			log.Println("send command", command.Args)
			fmt.Fprintf(conn, command.Args + "\n")	
			message, _ := bufio.NewReader(conn).ReadString('\n')
			return message, nil	
		}	else {
			return "disconnected", nil
		}
	}

	return "", nil

}
