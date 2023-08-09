package message

import (
	"encoding/json"
)

type TextMessage struct {
	BaseMessage
	Text string `json:"text"`
}

func (t *TextMessage) GetMessageType() string {
	return "text"
}
func (t *TextMessage) ToJson() string {
	jsonString, _ := json.Marshal(t)
	return string(jsonString)
}
func (t *TextMessage) FromJson(data json.RawMessage) {
	json.Unmarshal([]byte(data), &t)
}
