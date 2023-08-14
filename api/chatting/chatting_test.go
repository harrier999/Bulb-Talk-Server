package chatting

import (
	"encoding/json"
	"log"
	"server/internal/models/message"
	"testing"
)

func TestJsonUnmarshal(t *testing.T) {
	var msg message.Message
	json.Unmarshal([]byte(`{"type":"message","data":"hello"}`), &msg)
	log.Println(msg)
}
