package zpclient

import (
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type WsClient struct {
	id     int
	ws     *websocket.Conn
	ch     chan string
	RcvCh  chan string
	doneCh chan bool
}

// Create a new connection
func NewWsClient(url string) *WsClient {

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	log.Println("Create new websocket connection")
	if err != nil {
		log.Fatal("dial:", err)
	}
	maxId++
	ch := make(chan string, channelBufSize)
	RcvCh := make(chan string, channelBufSize)
	doneCh := make(chan bool)

	return &WsClient{maxId, ws, ch, RcvCh, doneCh}
}

func (c *WsClient) Conn() *websocket.Conn {
	return c.ws
}

func (c *WsClient) Write(msg string) {
	select {
	case c.ch <- msg:
	default:
		fmt.Errorf("client %d is disconnected.", c.id)
	}
}

func (c *WsClient) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *WsClient) Listen() {
	go c.listenWrite()
	c.listenRead()
	log.Println("End listenRead");
}

// Listen write request via chanel
func (c *WsClient) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			err := c.ws.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("Send message error %#v\n", err)
			}

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *WsClient) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			_, message, err := c.ws.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("recv: %s", message)
			c.RcvCh <- string(message)
			if err == io.EOF {
				c.doneCh <- true
			}
		}
	}
}
