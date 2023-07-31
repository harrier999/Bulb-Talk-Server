package chatting

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
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

func Handler(w http.ResponseWriter, r *http.Request) {
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

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
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

type ChattingScoket struct {

	

}

type userInfo struct {
	gorm.Model
	Username string `json:"username"`

}
