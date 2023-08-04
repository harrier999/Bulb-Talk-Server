package chatting

import (
	"context"
	"log"
	"net/http"

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

func GetChattingHistory(room_id int, last_message_id int) {
	return 

}

type User struct {
	Id string `json:"id"`
	UserName string `json:"firstName"`
	ProfileImage string `json:"imageUrl"`
}

type Message struct {
	Id string `json:"id"`
	RoomId string `json:"roomId"`
	Type string `json:"type"`
	Author User `json:"author"`
}
type Messenger interface {
	GetMessageType() string
	ToJson() string
	FromJson(string)
}
type TextMessage struct {
	Message
	Text string `json:"text"`
}
// type ImageMessage struct {
// 	Message
// 	ImageUrl string `json:"uri"`
// 	Name string `json:"name"`
// 	Size int `json:"size"`
// }
// type Request struct{
// 	RequestType string `json:"requestType"`
// 	Playload map[string]json.RawMessage `json:"playload"`
// } 

// type Wrapper interface {
// 	GetRequestType() string
// 	FromJson(string)
// 	CreateResponse() string
// }
// type Response{

// }
	

// type ChattingScoket struct {

	

// }

// type userInfo struct {
// 	gorm.Model
// 	Username string `json:"username"`

// }
