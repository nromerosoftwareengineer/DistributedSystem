package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

var app_context *AppContext

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

type AppContext struct {
	userId_websocket_map map[string]*websocket.Conn
	redisClient          *redis.Client
	ctx                  context.Context
}

func init_app_context() {
	app_context = &AppContext{
		userId_websocket_map: make(map[string]*websocket.Conn),
		ctx:                  context.Background(),
	}
	init_redis(app_context)
}

func (wsh webSocketHandler) upgrade_to_ws(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err == nil {
		return c
	}
	log.Printf("error %s when upgrading connection to websocket", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write([]byte("Unable to upgrade to websocket"))
	if err != nil {
		log.Printf("error %s when trying to return error message to client", err)
	}
	return nil
}

func handle_if_client_closed_connection(userId string, err error, c *websocket.Conn) {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		err := c.Close()
		if err != nil {
			log.Printf("error %s when trying to close websocket connectin for userId:%s", userId, err)
			return
		}
		log.Println("Connection closed by client")
	} else {
		log.Println("Error reading JSON:", err)
	}
}

func message_loop_handler(userId string) {
	var msg Message
	c := app_context.userId_websocket_map[userId]
	for {
		err := c.ReadJSON(&msg)
		if err != nil {
			handle_if_client_closed_connection(userId, err, c)
			break
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %v\n", err)
			continue
		}
		err = app_context.redisClient.Publish(app_context.ctx, "chat-message", data).Err()
		log.Printf("Received message from userId:%s to userId:%s \n", msg.From, msg.To)
		if err != nil {
			log.Printf("Error publishing message: %v\n", err)
			continue
		}
	}
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	log.Printf("userId: %s is trying to connect, upgrading to websocket connection\n", userId)
	c := wsh.upgrade_to_ws(w, r)
	if c == nil {
		return
	}
	app_context.userId_websocket_map[userId] = c
	message_loop_handler(userId)
}

func main() {
	init_app_context()
	go startSubscriber(context.Background(), app_context.redisClient, "chat-message")
	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8100", nil))
}

func init_redis(appContext *AppContext) {
	run_on_container := os.Getenv("IS_CONTAINER_RUN")
	var address string
	if run_on_container == "true" {
		address = "redis:6379"
		log.Println("Redis running on " + address)
	} else {
		log.Println("Redis running on " + address)
		address = "0.0.0.0:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	pong, err := rdb.Ping(appContext.ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis connected:", pong)
	appContext.redisClient = rdb
}

func startSubscriber(ctx context.Context, rdb *redis.Client, channel string) {
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	fmt.Printf("Subscribed to channel: %s\n", channel)
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatalf("Error receiving message: %v\n", err)
		}

		var message Message
		err = json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			log.Printf("Invalid JSON message: %s\n", msg.Payload)
			continue
		}

		v := app_context.userId_websocket_map[message.To]

		if v == nil {
			continue
		}

		writeMessage := Message{
			From: message.From,
			To:   message.To,
			Body: message.Body,
		}
		err = v.WriteJSON(writeMessage)
		if err != nil {
			log.Println("Error writing JSON:", err)
		}
	}
}
