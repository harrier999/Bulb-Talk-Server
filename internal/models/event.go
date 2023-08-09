package models

import (
	"encoding/json"
	"server/internal/models/message"
)

type Event struct {
	MessageType string          `json:"messageType"`
	Payload     json.RawMessage `json:"payload"`
}

func jsonToMessage(jsonString string) Event {
	var event Event
	json.Unmarshal([]byte(jsonString), &event)
	return event
}
func MessageToJson(message message.Message) string {
	jsonString, _ := json.Marshal(message)
	return string(jsonString)
}
