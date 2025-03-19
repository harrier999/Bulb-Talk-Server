package message

import (
	"encoding/json"
	"time"
)

type TextMessage struct {
	BaseMessage
	Content string `json:"content"`
}

func (t *TextMessage) GetMessageType() string {
	return "message"
}

func (t *TextMessage) ToJson() string {
	if t.Timestamp == "" {
		t.Timestamp = time.Now().Format(time.RFC3339)
	}
	jsonString, _ := json.Marshal(t)
	return string(jsonString)
}

func (t *TextMessage) FromJson(data json.RawMessage) {
	json.Unmarshal([]byte(data), &t)
}
