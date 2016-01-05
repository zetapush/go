package zpclient

import (
	"encoding/json"
)

type Message struct {
	ClientId                 string           `json:"clientId"`
	Data                     *json.RawMessage `json:"data,omitempty"`
	Channel                  string           `json:"channel,omitempty"`
	Id                       int              `json:"id"`
	Error                    string           `json:"error,omitempty"`
	Timestamp                string           `json:"timestamp,omitempty"`
	Transport                string           `json:"transport,omitempty"`
	Advice                   *Advice          `json:"advice,omitempty"`
	Successful               bool             `json:"successful,omitempty"`
	Subscription             string           `json:"subscription,omitempty"`
	Extension                *json.RawMessage `json:"ext,omitempty"`
	ConnectionType           string           `json:"connectionType,omitempty"`
	Version                  string           `json:"version,omitempty"`
	MinimumVersion           string           `json:"minimumVersion,omitempty"`
	SupportedConnectionTypes []string         `json:"supportedConnectionTypes,omitempty"`
	AuthSuccessful           bool             `json:"authSuccessful,omitempty"`
	wait                     bool
}

type Advice struct {
	Reconnect   string `json:"reconnect,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`
	Interval    int    `json:"interval,omitempty"`
	MaxInterval int    `json:"maxInterval,omitempty"`
}

type Config struct {
	StickyReconnect           bool
	ConnectTimeout            int
	MaxConnections            int
	BackoffIncrement          int
	MaxBackoff                int
	LogLevel                  string
	ReverseIncomingExtensions bool
	MaxNetworkDelay           int
	AppendMessageTypeToURL    bool
	AutoBatch                 bool
	MaxURILength              int
	Advice                    Advice
}
