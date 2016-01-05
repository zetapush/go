package zpservice

import (
	"encoding/json"
)

type ErrorMessage struct {
	Message     string           `json:"message"`
	Code        string           `json:"code"`
	ChannelName string           `json:"channelName"`
	Source      *json.RawMessage `json:"source"`
}
