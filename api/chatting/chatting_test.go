package chatting

import (
	"testing"
	"log"
	"encoding/json"
)

var request = []byte(`{"requestType":"message", "payload":"hello"}`)

func TestJsonUnmarshal(t *testing.T){
	var msg Message
	json.Unmarshal([]byte(`{"type":"message","data":"hello"}`), &msg)
	log.Println(msg)
}