package zpservice

import (
	"encoding/json"
)

type MacroPlay struct {
	Debug    int              `json:"debug,omitempty"`
	HardFail bool             `json:"hardFail,omitempty"`
	Name     string           `json:"name"`
	Params   *json.RawMessage `json:"params,omitempty"`
}

type MacroCompletion struct {
	ElapsedMillis int64                 `json:"elapsedMillis`
	Name          string                `json:"name"`
	NbCalls       int64                 `json:"nbCalls`
	Result        *json.RawMessage      `json:"result"`
	Log           []string              `json:"log"`
	Errors        []MacroExecutionError `json:"errors"`
}

type MacroExecutionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Step    string `json:"step"`
}
