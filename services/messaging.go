package zpservice

import (
	"encoding/json"
)

type MessagingMessage struct {
	Target  string           `json:"target"`
	Channel string           `json:"channel,omitempty"`
	Data    *json.RawMessage `json:"data"`
}

type MessagingResult struct {
	Target  string           `json:"target"`
	Source  string           `json:"source"`
	Channel string           `json:"channel,omitempty"`
	Data    *json.RawMessage `json:"data"`
}
