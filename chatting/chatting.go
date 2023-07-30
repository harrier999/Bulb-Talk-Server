package chatting

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections = make(map[string]*websocket.Conn)

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Println(err)
		return
	}
	user_id := r.Header.Get("user_id")
	connections[user_id] = conn
	defer delete (connections, user_id)
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		for user, conn := range connections {
			err = conn.WriteMessage(websocket.TextMessage, msg)
			fmt.Println("user: ", user);
			fmt.Println(string(msg))
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