# Info
Go client for ZetaPush

# How to use ?

## Import Client & Services
import (
  "encoding/json"
	"github.com/zetapush/go/client"
	"github.com/zetapush/go/services"
	)

## Initialize a zpclient object

zpC := zpclient.Client{}
zpC.Init("YourBusinessId")

## Create and use services

// Create a service with his deploymentId
echoService := zpclient.CreateService(&zpC, "deploymentId")

// Add a callback for the specific verb of the service
echoService.On("echo", func(m *zpclient.Message) {
		log.Printf("received a message from echoService %#v\n", m)

		var echoMessage zpservice.EchoMessage
    // Unmarshal the message data you've just received
		if err := json.Unmarshal([]byte(*m.Data), &echoMessage); err != nil {
			log.Fatal(err)
		}

		fmt.Println("echoMessage", echoMessage.Message)

	})
	
## Connect to ZetaPush with an authentication

	simpleAuthent := zpclient.NewSimpleAuthentication(&zpC, "deploymentId")
	simpleAuthent.Login = "login"
	simpleAuthent.Password = "pwd"
	simpleAuthent.Resource = "resource"
	// A callback when you're connected
	simpleAuthent.OnConnected(func() {
		log.Println("OnConnected ", simpleAuthent.UserId, simpleAuthent.RMToken)
			})
	// Call the connect method of your zpC client		
	zpC.Connect(simpleAuthent)
			
