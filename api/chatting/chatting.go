package chatting

import (
	"context"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/websocket"

	//"github.com/jinzhu/gorm"
	"server/internal/db"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func closeConnection(room_id string, user_id string){
	connections[room_id][user_id].Close()
	delete(connections[room_id], user_id)
	log.Println("Closing connection from user ", user_id, "room ", room_id)
	if len(connections[room_id]) == 0 {
		delete(connections, room_id)
	}
	
}


var connections = make(map[string]map[string]*websocket.Conn)
var ctx context.Context
var redisDB *redis.Client

// func Init() {
// 	ctx = db.GetContext()
// 	redisDB = db.GetRedisClient()
// }

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx = db.GetContext()
	redisDB = db.GetRedisClient()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Println(err)
		return
	}
	room_id := r.Header.Get("room_id")
	user_id := r.Header.Get("user_id")
	log.Println("new connection from ", user_id)
	log.Println("room_id is ", room_id)

	if connections[room_id] == nil{
		connections[room_id] = make(map[string]*websocket.Conn)
		log.Println("creating new room")
	}
	
	connections[room_id][user_id] = conn
	defer closeConnection(room_id, user_id)

	chattingHistory := redisDB.LRange(ctx, "room_id:"+room_id, 0, -1)


	for _, msg := range chattingHistory.Val() {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}





	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		redisDB.RPush(ctx, "room_id:"+room_id, string(msg))
		log.Println("New Message!")
		log.Println(string(msg))
		for user, conn := range connections[room_id] {
			log.Println("delevering to user ", user)
			
			err = conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
		
	}

}

func GetChattingHistory(room_id string, last_message_id int64, redisDB *redis.Client) ([]string, error) {
	chattingHistory := redisDB.LRange(ctx, "room_id:"+string(room_id), last_message_id, -1)

	if chattingHistory.Err() != nil {
		log.Println(chattingHistory.Err())
		return nil, chattingHistory.Err()
	}
	return chattingHistory.Val(), nil
}

func redisListToMessageList(redisList []string) []Chat {
	var messageList []Chat
	var chat map[string]interface{}

	for _, msg := range redisList {
		err := json.Unmarshal([]byte(msg), &chat)
		if err != nil{
			return messageList
			// TODO
		}
		switch chat["type"] {
		case "text":
			json.Marshal(msg)
			
		}

	}
	return messageList
}


type User struct {
	Id string `json:"id"` 
	// UserName string `json:"firstName"`
	// ProfileImage string `json:"imageUrl"`
}
type Chat interface {
	GetMessageType() string
	ToJson() string
	FromJson(string)
}
type BaseChat struct {
	Id string `json:"id"`
	RoomId string `json:"roomId"`
	Type string `json:"type"`
	Author User `json:"author"`
}
type TextChat struct {
	BaseChat
	Text string `json:"text"`
}
type Message struct{
	Message string `json:"MessageType"`
	Payload json.RawMessage `json:"payload"`
} 
// type ImageChat struct {
// 	BaseChat
// 	ImageUrl string `json:"uri"`
// 	Name string `json:"name"`
// 	Size int `json:"size"`
// }


func jsonToMessage (jsonString string) Message{
	var message Message
	json.Unmarshal([]byte(jsonString), &message)
	return message
}
func MessageToJson(message Message) string{
	jsonString, _ := json.Marshal(message)
	return string(jsonString)
}
func (t *TextChat) GetMessageType() string{
	return "text"
}
func (t *TextChat) ToJson() string{
	jsonString, _ := json.Marshal(t)
	return string(jsonString)
}
func (t *TextChat) FromJson(data json.RawMessage){
	json.Unmarshal([]byte(data), &t)
}




