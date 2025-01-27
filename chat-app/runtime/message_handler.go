package runtime

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

const ChatMessageChannel = "chat-message"

type MessageHandler struct {
	ch          *ConnectionHandler
	RedisClient *Redis
}

func NewMessageHandler(handler *ConnectionHandler, redisClient *Redis) *MessageHandler {
	mh := &MessageHandler{
		ch:          handler,
		RedisClient: redisClient,
	}
	return mh
}

func (mh *MessageHandler) HandleMessageLoop(userId string) {
	var msg Message
	c := mh.ch.GetUserConn(userId)
	for c != nil {
		err := c.ReadJSON(&msg)

		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			mh.ch.CloseUserConn(userId)
			break
		}

		if err != nil {
			log.Printf("websocket connection closed %v\n", err)
			mh.ch.CloseUserConn(userId)
			break
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %v\n", err)
			continue
		}
		log.Printf("Received message from userId:%s to userId:%s \n", msg.From, msg.To)
		mh.RedisClient.Publish(data, ChatMessageChannel)
		c = mh.ch.GetUserConn(userId)
	}
}

func (mh *MessageHandler) HandlePublishedMessages() {
	mh.RedisClient.HandlePublishedMessages(mh)
}

func (mh *MessageHandler) Close() {
	mh.RedisClient.Close()
}

func (mh *MessageHandler) SendMessageToUsers(message *Message, users []string) {
	for _, user := range users {
		mh.sendMessageTo(message, user)
	}
}

func (mh *MessageHandler) sendMessageTo(message *Message, to string) {

	c := mh.ch.GetUserConn(to)

	if c == nil {
		log.Printf("WebSocket connection not found for user: %s", message.To)
		return
	} else {
		log.Printf("Found active WebSocket connection for user: %s", message.To)
	}

	writeMessage := message
	writeMessage.To = to

	err := c.WriteJSON(writeMessage)
	if err != nil {
		log.Println("Error writing message:", err)
	}
}
