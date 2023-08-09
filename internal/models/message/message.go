package message

import (
	"encoding/json"
)

type User struct {
	Id string `json:"id"`
}

type Message interface {
	GetMessageType() string
	ToJson() string
	FromJson(data json.RawMessage)
}
type BaseMessage struct {
	Id     string `json:"id"`
	RoomId string `json:"roomId"`
	Type   string `json:"type"`
	Author User   `json:"author"`
}

type MessageList struct {
	FirstSeq int64     `json:"firstSeq"`
	LastSeq  int64     `json:"lastSeq"`
	Messages []Message `json:"messages"`
}
