package message

import (
	"encoding/json"

	"github.com/google/uuid"
)

type User struct {
	Id string `json:"id"`
}

type Message interface {
	GetID() uuid.UUID
	GetRoomID() string
	GetType() string
	GetAuthor() User
	GetTimestamp() string
	GetMessageType() string
	ToJson() string
	FromJson(data json.RawMessage)
}

type BaseMessage struct {
	Id        uuid.UUID `json:"id"`
	RoomId    string    `json:"roomId"`
	Type      string    `json:"type"`
	Author    User      `json:"author"`
	Timestamp string    `json:"timestamp,omitempty"`
}

func (m *BaseMessage) GenerateID() {
	id, _ := uuid.NewV7()
	m.Id = id
}

func (m *BaseMessage) GetID() uuid.UUID {
	return m.Id
}

func (m *BaseMessage) GetRoomID() string {
	return m.RoomId
}

func (m *BaseMessage) GetType() string {
	return m.Type
}

func (m *BaseMessage) GetAuthor() User {
	return m.Author
}

func (m *BaseMessage) GetTimestamp() string {
	return m.Timestamp
}

func (m *BaseMessage) GetMessageType() string {
	return m.Type
}

func (m *BaseMessage) ToJson() string {
	jsonData, _ := json.Marshal(m)
	return string(jsonData)
}

func (m *BaseMessage) FromJson(data json.RawMessage) {
	json.Unmarshal(data, m)
}

type MessageResponse struct {
	Success  bool      `json:"success"`
	Messages []Message `json:"messages"`
}
