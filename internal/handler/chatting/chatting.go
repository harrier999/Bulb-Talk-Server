package chatting

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"server/internal/models"
	"server/internal/models/message"

	"github.com/gorilla/websocket"
	"server/internal/db/redis_db"

	"github.com/redis/go-redis/v9"
)

var ALLOWED_ORIGIN_LIST = []string{"app.bulbtalk.com"}

type initialMessage struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		url, err := url.Parse(origin)
		if err != nil {
			log.Println(err)
			return false
		}
		for _, allowed_origin := range ALLOWED_ORIGIN_LIST {
			if url.Host == allowed_origin {
				return true
			}
		}
		return false
	},
}

func closeConnection(room_id string, user_id string) {

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
	log.Println("new connection")
	ctx = redis_db.GetContext()
	redisDB = redis_db.GetChattingHistoryClient()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	// room_id := r.Header.Get("room_id")
	// user_id := r.Header.Get("user_id")

	var initialMessage initialMessage
	err = json.Unmarshal(msg, &initialMessage)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	room_id := initialMessage.RoomID
	user_id := initialMessage.UserID

	log.Println("new connection from ", user_id)
	log.Println("room_id is ", room_id)

	if connections[room_id] == nil {
		connections[room_id] = make(map[string]*websocket.Conn)
		log.Println("creating new room")
	}

	connections[room_id][user_id] = conn
	defer closeConnection(room_id, user_id)

	chattingHistory, err := GetChattingHistory(room_id, 0, redisDB)
	if err != nil {
		log.Println(err)
		log.Println("error in getting chatting history")
		return
	}

	conn.WriteMessage(websocket.TextMessage, chattingHistory)
	log.Println(string(chattingHistory))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(msg))
		redisDB.RPush(ctx, "room_id:"+room_id, string(msg))
		log.Println("New Message!")
		log.Println(string(msg))
		eventData := models.Event{MessageType: "textMessage", Payload: msg}
		event, err := json.Marshal(eventData)
		if err != nil {
			log.Println(err)
			return
		}
		for user, conn := range connections[room_id] {
			log.Println("delevering to user ", user)

			err = conn.WriteMessage(websocket.TextMessage, event)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func GetChattingHistoryFromRedis(room_id string, last_message_id int64, redisDB *redis.Client) ([]string, error) {
	chattingHistory := redisDB.LRange(ctx, "room_id:"+string(room_id), last_message_id, -1)

	if chattingHistory.Err() != nil {
		log.Println(chattingHistory.Err())
		return nil, chattingHistory.Err()
	}
	return chattingHistory.Val(), nil
}

func redisListToMessageList(redisList []string) ([]message.Message, error) {
	var messageList []message.Message
	var singleMessage message.Message
	type typer struct {
		Type string `json:"type"`
	}
	var t typer

	for _, msg := range redisList {
		err := json.Unmarshal([]byte(msg), &t)
		if err != nil {
			return nil, err
		}

		switch t.Type {
		case "text":
			singleMessage = &message.TextMessage{}
		case "image":
			singleMessage = &message.ImageMessage{}
		default:
			log.Println("error in getting message type")
			log.Println("error is ", t.Type)
			return nil, errors.New("error in getting message type")
		}
		singleMessage.FromJson([]byte(msg))
		messageList = append(messageList, singleMessage)
		log.Println(singleMessage)

	}
	return messageList, nil
}

func GetChattingHistory(room_id string, last_message_id int64, redisDB *redis.Client) (json.RawMessage, error) {
	chattingHistoryRaw, err := GetChattingHistoryFromRedis(room_id, last_message_id, redisDB)
	if err != nil {
		log.Println(err)
		log.Println("error in getting chatting history from redis")
		return nil, err
	}
	chattingHistory, err := redisListToMessageList(chattingHistoryRaw)
	if err != nil {
		log.Println(err)
		log.Println("error in converting redis list to message list")
		return nil, err
	}
	messageListWithSeq := message.MessageList{FirstSeq: 0, LastSeq: int64(len(chattingHistory)), Messages: chattingHistory}
	jsonString, err := json.Marshal(messageListWithSeq)
	if err != nil {
		log.Println(err)
		log.Println("error in converting message list to json")
		return nil, err
	}
	event := models.Event{MessageType: "messageList", Payload: jsonString}
	json, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		log.Println("error in converting event to json")
		return nil, err
	}
	return json, nil

}
