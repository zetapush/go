package zpclient

import (
	"encoding/json"
)

type Service struct {
	BusinessId   string
	DeploymentId string
	client       *Client
}

type MsgHandler func(message *Message)

type Subscription struct {
	Channel string
	closed  bool
	mcb     MsgHandler
	mch     chan interface{}
	sc      bool
}

func (s *Service) getChannel(verb string) string {
	return "/service/" + s.BusinessId + "/" + s.DeploymentId + "/" + verb
}

func (s *Service) Send(verb string, data interface{}) {

	message := &Message{Channel: s.getChannel(verb)}

	returnByteArray, _ := json.Marshal(data)
	jsonRaw := json.RawMessage(returnByteArray)
	message.Data = &jsonRaw

	s.client.sendMessage(message)
}

func (s *Service) On(verb string, callback MsgHandler) *Subscription {

	topic := s.getChannel(verb)

	sub := &Subscription{Channel: topic, mcb: callback}

	sub.mch = s.client.ps.Sub(topic)

	if callback != nil {
		go s.waitMessage(sub)

	}

	return sub
}

func (s *Service) waitMessage(sub *Subscription) {
	mcb := sub.mcb
	ch := sub.mch

	if ch == nil {
		return
	}
	for {
		m, ok := <-ch
		if !ok {
			break
		}
		// call the callback
		receiveMsg := m.(*Message)
		mcb(receiveMsg)
	}
}

func CreateService(client *Client, deploymentId string) *Service {

	s := &Service{BusinessId: client.BusinessId, DeploymentId: deploymentId}
	s.client = client

	return s
}
