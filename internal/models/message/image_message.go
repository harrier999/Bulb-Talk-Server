package message

import (
	"encoding/json"
	"time"
)

// ImageMessage는 이미지 메시지를 나타냅니다.
type ImageMessage struct {
	BaseMessage
	ImageURL string `json:"imageUrl"`
}

func (i *ImageMessage) GetMessageType() string {
	return "image"
}

func (i *ImageMessage) ToJson() string {
	if i.Timestamp == "" {
		i.Timestamp = time.Now().Format(time.RFC3339)
	}
	jsonString, _ := json.Marshal(i)
	return string(jsonString)
}

func (i *ImageMessage) FromJson(data json.RawMessage) {
	json.Unmarshal([]byte(data), &i)
}
