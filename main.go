package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bennyscetbun/jsongo"
	"log"
	"os"
	"zetapush/client"
	"zetapush/services"
)

func main() {

	zpC := zpclient.Client{}
	zpC.Init("GmY-HuzW")

	echoService := zpclient.CreateService(&zpC, "w3FQ")

	echoService.On("echo", func(m *zpclient.Message) {
		log.Printf("received a message from echoService %#v\n", m)

		var echoMessage zpservice.EchoMessage

		if err := json.Unmarshal([]byte(*m.Data), &echoMessage); err != nil {
			log.Fatal(err)
		}

		fmt.Println("echoMessage", echoMessage.Message)

	})

	macroService := zpclient.CreateService(&zpC, "57C3")

	macroService.On("error", func(m *zpclient.Message) {

		log.Printf("received an error message from macroService %#v\n", m)

		var errorMessage zpservice.ErrorMessage

		if err := json.Unmarshal([]byte(*m.Data), &errorMessage); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Error Message %#v\n", errorMessage)
	})

	macroService.On("completed", func(m *zpclient.Message) {

		log.Printf("received a message from macroService %#v\n", m)

		var macroCompletion zpservice.MacroCompletion

		if err := json.Unmarshal([]byte(*m.Data), &macroCompletion); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("macroCompletion %#v\n", macroCompletion)
	})

	gdaService := zpclient.CreateService(&zpC, "IFa0")

	gdaService.On("get", func(m *zpclient.Message) {

		log.Printf("received a message from gdaServce %#v\n", m)

		var gdaGetResult zpservice.GdaGetResult
		if err := json.Unmarshal([]byte(*m.Data), &gdaGetResult); err != nil {
			log.Fatal(err)
		}

		root := jsongo.JSONNode{}

		err := json.Unmarshal([]byte(*gdaGetResult.Result), &root)

		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		root.DebugProspect(0, "\t")
		// Get the index __key value
		log.Println(root.At("__key").Get().(string))

	})

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		key := string([]byte(input)[0])
		switch key {
		case "x":
			fmt.Println("type x")
		case "q":
			zpC.Disconnect()
			return
		case "c":
			// Connect
			simpleAuthent := zpclient.NewSimpleAuthentication(&zpC, "KZyH")
			simpleAuthent.Login = "tuto"
			simpleAuthent.Password = "tuto"
			simpleAuthent.Resource = "testGo"
			simpleAuthent.OnConnected(func() {
				log.Println("OnConnected ", simpleAuthent.UserId, simpleAuthent.RMToken)

			})
			zpC.Connect(simpleAuthent)
		case "m":
			macroPlayMessage := zpservice.MacroPlay{Debug: 4, Name: "chatroom"}

			macroService.Send("call", macroPlayMessage)
		case "g":
			gdaGetMessage := zpservice.GdaGet{Table: "userList", Key: "123"}

			gdaService.Send("get", gdaGetMessage)
		case "t":
			echoMessage := zpservice.EchoMessage{Message: "hello you!"}

			echoService.Send("echo", echoMessage)

		}

	}
}
