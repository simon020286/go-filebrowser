package ws

import (
	"encoding/json"
)

// Message struct definition.
type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ToByte serialize object.
func (msg *Message) ToByte() ([]byte, error) {
	return json.Marshal(*msg)
}
