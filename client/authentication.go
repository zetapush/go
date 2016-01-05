package zpclient

import (
	"encoding/json"
	"log"
)

type loginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginTokenData struct {
	Token string `json:"token"`
}

type MetaAuthenticationResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type ExtAuthentResponse struct {
	Authentication MetaAuthenticationResponse `json:"authentication,omitempty"`
}

type MetaAuthentication struct {
	Action   string      `json:"action"`
	Type     string      `json:"type"`
	Resource string      `json:"resource"`
	Data     interface{} `json:"data"`
}

type ExtAuthent struct {
	Authentication MetaAuthentication `json:"authentication"`
}

type Authentication struct {
	client            *Client
	UserId            string
	ClientId          string
	Token             string
	BusinessId        string
	DeploymentId      string
	Resource          string
	RMToken           string // rememberMe
	authentMsgHandler AuthentMsgHandler
}

type AuthentMsgHandler func()

type AuthenticationInterface interface {
	GetConnectionData() json.RawMessage
	GetClientId() string
	OnConnected(callback AuthentMsgHandler)
}

type Simple struct {
	Authentication
	Login    string
	Password string
}

type Weak struct {
	Authentication
}

func NewSimpleAuthentication(client *Client, deploymentId string) *Simple {
	simple := &Simple{}
	simple.client = client
	simple.BusinessId = client.BusinessId
	simple.DeploymentId = deploymentId
	return simple
}

/*
	Create an extension authent with login/pwd
*/
func (s *Simple) GetConnectionData() json.RawMessage {

	authentication := MetaAuthentication{Action: "authenticate", Resource: s.Resource}
	authentication.Type = s.BusinessId + "." + s.DeploymentId + ".simple"
	data := loginData{Login: s.Login, Password: s.Password}
	authentication.Data = data

	extAuthent := ExtAuthent{Authentication: authentication}
	returnByteArray, _ := json.Marshal(extAuthent)

	go s.waitForHandshake()

	return json.RawMessage(returnByteArray)
}

func (s *Simple) OnConnected(callback AuthentMsgHandler) {
	s.authentMsgHandler = callback
}

func (s *Simple) waitForHandshake() {
	ch := s.client.ps.SubOnce(META_HANDSHAKE)

	m, ok := <-ch
	if !ok {
		return
	}

	receiveMsg := m.(*Message)
	var extAuthentResponse ExtAuthentResponse
	if err := json.Unmarshal([]byte(*receiveMsg.Extension), &extAuthentResponse); err != nil {
		log.Fatal(err)
	}

	s.UserId = extAuthentResponse.Authentication.UserId
	s.RMToken = extAuthentResponse.Authentication.Token

	if s.authentMsgHandler != nil {
		s.authentMsgHandler()
	}
	log.Println("waitForHanshake OK", s.UserId, s.RMToken)
}

func (s *Simple) GetClientId() string {
	return s.ClientId
}

func (w *Weak) GetConnectionData() json.RawMessage {
	return nil
}

func (w *Weak) GetClientId() string {
	return ""
}
func GetConnectionData(authenticationInterface AuthenticationInterface) json.RawMessage {
	return authenticationInterface.GetConnectionData()
}
