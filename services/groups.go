package zpservice

import (
	"encoding/json"
)

const (
	BAD_INPUT_DATA = "Missing group id."
	BAD_NAME       = "Group id must be alphanumerical."
)

/*Remoting*/
type RemoteCommand struct {
	Cmd          string           `json:"cmd"`
	From         string           `json:"from,omitempty"`
	FromResource string           `json:"fromResource,omitempty"`
	Owner        string           `json:"owner,omitempty"`
	Resource     string           `json:"resource,omitempty"`
	Data         *json.RawMessage `json:"data,omitempty"`
}

type DeviceCapabilities struct {
	AnsweringResource string   `json:"answeringResource,omitempty"`
	AskingResource    string   `json:"askingResource,omitempty"`
	Capabilities      []string `json:"capabilities,omitempty"`
}

type PingRequest struct {
	Action    string `json:"action"`
	Available bool   `json:"available,omitempty"`
	Owner     string `json:"owner,omitempty"`
	User      string `json:"user,omitempty"`
}

type DeviceAvailability struct {
	Action    string `json:"action,omitempty"`
	Available bool   `json:"available,omitempty"`
	Uid       string `json:"uid,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Resource  string `json:"resource,omitempty"`
	User      string `json:"user,omitempty"`
}

/*Groups and rights*/
type UserGroup struct {
	Group string `json:"group"`
	Owner string `json:"owner,omitempty"`
	User  string `json:"user,omitempty"`
}

type GroupUsers struct {
	Group string   `json:"group"`
	Owner string   `json:"owner,omitempty"`
	Users []string `json:"users,omitempty"`
}

type GroupInfo struct {
	Group     string `json:"group"`
	GroupName string `json:"groupName"`
	Owner     string `json:"owner,omitempty"`
}

type GroupRelated struct {
	Group string `json:"group"`
	Owner string `json:"owner,omitempty"`
}

type Grant struct {
	Action   string `json:"action,omitempty"`
	Group    string `json:"group,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type Grants struct {
	Actions  string `json:"actions,omitempty"`
	Group    string `json:"group,omitempty"`
	Owner    string `json:"owner,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type GrantList struct {
	Grants []GrantListItem `json:"grants"`
	Group  string          `json:"group,omitempty"`
	Owner  string          `json:"owner,omitempty"`
}

type GrantListItem struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type GroupPresence struct {
	Group     string     `json:"group,omitempty"`
	Owner     string     `json:"owner,omitempty"`
	Presences []Presence `json:"presences,omitempty"`
}

type OwnerResource struct {
	Owner    string `json:"owner,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type Presence struct {
	Group    GroupRelated  `json:"group,omitempty"`
	Presence string        `json:"presence,omitempty"`
	User     OwnerResource `json:"user,omitempty"`
}

type UserMembership struct {
	Group    string `json:"group"`
	Owner    string `json:"owner,omitempty"`
	HardFail bool   `json:"hardFail,omitempty"`
}

type UserGroupMembership struct {
	Group  string `json:"group,omitempty"`
	Member string `json:"member,omitempty"`
	User   string `json:"user,omitempty"`
	Owner  string `json:"owner,omitempty"`
}
