package chatting

import (
	"context"
	"log"
	"net/http"
	"encoding/json"

	"server/internal/models/message"

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

func redisListToMessageList(redisList []string) []message.Message {
	var messageList []message.Message
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






