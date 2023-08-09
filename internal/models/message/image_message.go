package message

import (
	"encoding/json"
)

type ImageMessage struct {
	BaseMessage
	ImageUrl string `json:"uri"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
}

func (i *ImageMessage) GetMessageType() string {
	return "image"
}

func (i *ImageMessage) ToJson() string {
	jsonString, _ := json.Marshal(i)
	return string(jsonString)
}
func (i *ImageMessage) FromJson(data json.RawMessage) {
	json.Unmarshal([]byte(data), &i)
}
