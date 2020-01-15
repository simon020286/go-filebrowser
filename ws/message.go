package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Message struct definition.
type Message struct {
	ID    uuid.UUID   `json:"id"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// NewMessage message constructor.
func NewMessage(event string, data interface{}) Message {
	message := Message{ID: uuid.New(), Event: event, Data: data}
	return message
}

// ToByte serialize object.
func (msg *Message) ToByte() ([]byte, error) {
	return json.Marshal(*msg)
}
