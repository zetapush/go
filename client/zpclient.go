package zpclient

import (
	"encoding/json"
	"fmt"
	"github.com/tuxychandru/pubsub"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	META_HANDSHAKE    = "/meta/handshake"
	META_CONNECT      = "/meta/connect"
	META_UNSUCCESSFUL = "/meta/unsuccessful"
	META_DISCONNECT   = "/meta/disconnect"
	META_SUBSCRIBE    = "/meta/subscribe"
	META_UNSUBSCRIBE  = "/meta/unsubscribe"
	META_CONNECTED    = "/meta/connected"
)

var (
	wsc *WsClient

	clientId string = ""
)

type Client struct {
	connected                bool
	Timer                    int
	messageId                int
	doneCh                   chan bool
	_status                  string
	_clientId                string
	BusinessId               string
	ps                       *pubsub.PubSub
	timeout                  int
	serverList               ServerList
	_advice                  Advice
	_config                  Config
	_backoff                 int
	_scheduleSend            *time.Timer
	_authenticationInterface AuthenticationInterface
	_unconnectTime           time.Time
}

type ServerList struct {
	Servers []string `json:"servers"`
}

func (c *Client) getServer() string {
	if len(c.serverList.Servers) == 0 {
		log.Println("Call api.zpush.io")
		res, err := http.Get("http://api.zpush.io/" + c.BusinessId)
		if err != nil {
			// Todo: retry ??
			panic(err)
		}
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&c.serverList)
	}

	rand.Seed(42)

	return c.serverList.Servers[rand.Intn(len(c.serverList.Servers))]
}

func (c *Client) setStatus(newStatus string) {
	if c._status != newStatus {
		log.Println("Status", c._status, "->", newStatus)
		c._status = newStatus
	}
}

func (c *Client) isDisconnected() bool {
	return c._status == "disconnecting" || c._status == "disconnected"
}

func (c *Client) sendMessage(message *Message) {

	message.Id = c.messageId
	message.ClientId = c._clientId
	c.messageId++
	msg, _ := json.Marshal(message)
	wsc.Write("[" + string(msg) + "]")
}

func (c *Client) _updateAdvice(newAdvice Advice) {

	c._advice.Interval = newAdvice.Interval
	c._advice.Timeout = newAdvice.Timeout

	if len(newAdvice.Reconnect) > 0 {
		c._advice.Reconnect = newAdvice.Reconnect
	}
}

func (c *Client) Init(businessId string) {
	c.BusinessId = businessId
	c.serverList = ServerList{}
	c._config = Config{
		StickyReconnect:           true,
		ConnectTimeout:            0,
		MaxConnections:            2,
		BackoffIncrement:          1000,
		MaxBackoff:                60000,
		ReverseIncomingExtensions: true,
		MaxNetworkDelay:           10000,
		AppendMessageTypeToURL:    true,
		AutoBatch:                 false,
		MaxURILength:              2000,
		Advice:                    Advice{Reconnect: "retry", Timeout: 60000, Interval: 0},
	}

	c._updateAdvice(c._config.Advice)

	c.messageId = 1
	c._backoff = 0
	//url = "ws://m.zpush.ovh:8080/str/strd"
	c.doneCh = make(chan bool)
	c.ps = pubsub.New(1)
	u, err := url.Parse(c.getServer())
	_ = err
	serverUrl := "ws://" + u.Host + u.Path + "/strd"
	//fmt.Printf("url %#v\n", u, serverUrl)

	wsc = NewWsClient(serverUrl)
	go wsc.Listen()

	go c.listen()

}

/*
	Listen for incoming messages and MetaConnect
*/
func (c *Client) listen() {

	for {
		select {
		case <-c.doneCh:
			return
		case msg := <-wsc.RcvCh:

			mArray := make([]Message, 0)
			json.Unmarshal([]byte(msg), &mArray)
			for _, m := range mArray {
				c.manageIncoming(&m)
			}
		}
	}
	fmt.Println("listen finished ==============")
}

func (c *Client) manageIncoming(m *Message) {

	if m.Advice != nil {
		c._updateAdvice(*m.Advice)
	}
	switch {
	case m.Channel == META_HANDSHAKE:
		c.handshakeResponse(m)
	case m.Channel == META_CONNECT:
		c.connectResponse(m)
	case m.Channel == META_DISCONNECT:
		c.disconnectResponse(m)
	case m.Channel == META_SUBSCRIBE:
		c.subscribeResponse(m)
	case m.Channel == META_UNSUBSCRIBE:
		c.unsubscribeResponse(m)
	default:
		c.messageResponse(m)
	}

}

func (c *Client) handshakeResponse(m *Message) {

	log.Printf("handshakeResponse %#v\n", m)
	if m.Successful {
		c._clientId = m.ClientId

		c.ps.Pub(m, META_HANDSHAKE)

		action := "none"
		if !c.isDisconnected() {
			action = c._advice.Reconnect
		}
		switch action {
		case "retry":
			c._resetBackoff()
			c._delayedConnect()
		case "none":
			c._disconnect(true)
		default:
			log.Fatal("Unrecognized advice action " + action)
		}

	} else {
		c._failHandshake(m)
	}

}

func (c *Client) _failHandshake(m *Message) {

	log.Println("_failHandshake")
	c.ps.Pub(m, META_HANDSHAKE)
	c.ps.Pub(m, META_UNSUCCESSFUL)

	var retry = !c.isDisconnected() && c._advice.Reconnect != "none"
	if retry {
		c._increaseBackoff()
		c._delayedHanshake()
	} else {
		c._disconnect(true)
	}
}

func (c *Client) _cancelDelayedSend() {
	if c._scheduleSend != nil {
		c._scheduleSend.Stop()
	}
	c._scheduleSend = nil
}

func (c *Client) _delayedSend(f func()) {
	c._cancelDelayedSend()
	delay := c._advice.Interval + c._backoff
	log.Println("Function scheduled in ", delay, "ms, interval=", c._advice.Interval, "backoff =", c._backoff, f)
	c._scheduleSend = time.AfterFunc(time.Duration(delay)*time.Millisecond, f)
}

func (c *Client) _resetBackoff() {
	c._backoff = 0
}

func (c *Client) _increaseBackoff() {
	if c._backoff < c._config.MaxBackoff {
		c._backoff += c._config.BackoffIncrement
	}
}

func (c *Client) _delayedConnect() {
	c.setStatus("connecting")
	c._delayedSend(func() {
		c._connect()
	})
}

func (c *Client) _delayedHanshake() {
	c.setStatus("handshaking")
	c._delayedSend(func() {
		c._handshake()
	})
}

func (c *Client) _disconnect(abord bool) {
	c._cancelDelayedSend()

	c._clientId = ""
	c.setStatus("disconnected")
	c._resetBackoff()

}

func (c *Client) _failConnect(m *Message) {
	c.ps.Pub(m, META_CONNECT)
	c.ps.Pub(m, META_UNSUCCESSFUL)

	if c._unconnectTime.IsZero() {
		c._unconnectTime = time.Now()
	}

	action := "none"
	if !c.isDisconnected() {
		action = c._advice.Reconnect
	}

	maxInterval := c._advice.MaxInterval
	if maxInterval > 0 {
		expiration := c._advice.Timeout + c._advice.Interval + maxInterval
		unconnected := time.Since(c._unconnectTime)

		if unconnected.Nanoseconds()/int64(time.Millisecond)+int64(c._backoff) > int64(expiration) {
			action = "handshake"
		}
	}
	switch action {
	case "retry":
		c._delayedConnect()
		c._increaseBackoff()
	case "handshake":
		c._resetBackoff()
		c._delayedHanshake()
	case "none":
		c._disconnect(true)
	default:
		log.Fatal("Unrecognized advice action" + action)
	}

}

func (c *Client) _failDisconnect(m *Message) {
	c._disconnect(true)
	c.ps.Pub(m, META_DISCONNECT)
	c.ps.Pub(m, META_UNSUCCESSFUL)
}

func (c *Client) _handshake() {
	if c.isDisconnected() {
		c._updateAdvice(c._config.Advice)
	} else {
		c._advice.Reconnect = "retry"
	}

	// Get the extension part from the authentication used
	extAuthent := c._authenticationInterface.GetConnectionData()

	handshakeMessage := &Message{}
	handshakeMessage.Channel = META_HANDSHAKE
	handshakeMessage.SupportedConnectionTypes = []string{"websocket"}
	handshakeMessage.Extension = &extAuthent

	c.setStatus("handshaking")
	log.Printf("Handshake sent %#v\n", handshakeMessage)
	c.sendMessage(handshakeMessage)
}

func (c *Client) _connect() {

	if c.isDisconnected() {
		return
	}

	// Send a meta/connect message
	connectMessage := &Message{}
	connectMessage.Channel = META_CONNECT
	connectMessage.ConnectionType = "websocket"
	if !c.connected {
		advice := Advice{Timeout: 0}
		connectMessage.Advice = &advice
	}
	c.setStatus("connecting")
	c.sendMessage(connectMessage)
	c.setStatus("connected")
}

func (c *Client) connectResponse(m *Message) {

	wasConnected := c.connected
	c.connected = m.Successful
	if m.Successful {
		c.ps.Pub(m, META_CONNECT)

		if wasConnected {
			c.ps.Pub(m, META_CONNECTED)
		}
		action := "none"
		if !c.isDisconnected() {
			action = c._advice.Reconnect
		}

		switch action {
		case "retry":
			c._resetBackoff()
			c._delayedConnect()
		case "none":
			c._disconnect(false)
		default:
			log.Fatal("Unrecognized advice action " + action)
		}

	} else {
		log.Println("connectResponse - call _failConnect")
		c._failConnect(m)
	}
}

func (c *Client) disconnectResponse(m *Message) {
	if m.Successful {
		c._disconnect(false)
		c.ps.Pub(m, META_DISCONNECT)
	} else {
		c._failDisconnect(m)
	}
}

func (c *Client) subscribeResponse(m *Message) {

}

func (c *Client) unsubscribeResponse(m *Message) {

}

func (c *Client) messageResponse(m *Message) {
	log.Println("messageResponse")

	if m.Data != nil {
		c.ps.Pub(m, m.Channel)
	}
}

func (c *Client) Connect(authenticationInterface AuthenticationInterface) bool {

	c._authenticationInterface = authenticationInterface

	c._handshake()

	return true
}

/*
	Send a message
*/

func (c *Client) Disconnect() {
	c._disconnect(true)
	c.doneCh <- true
	wsc.Done()
}
