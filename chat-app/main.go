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
	"sync"
)

var app_context *AppContext

type webSocketHandler struct {
	upgrader websocket.Upgrader
	mapMutex sync.Mutex
}

func (app_context *AppContext) AddConnection(userId string, c *websocket.Conn) {
	app_context.mapMutex.Lock()
	defer app_context.mapMutex.Unlock()
	app_context.userId_websocket_map[userId] = c
}

func (app_context *AppContext) RemoveConnection(userId string) {
	app_context.mapMutex.Lock()
	defer app_context.mapMutex.Unlock()
	delete(app_context.userId_websocket_map, userId)
}

func (app_context *AppContext) GetConnection(userId string) *websocket.Conn {
	app_context.mapMutex.Lock()
	defer app_context.mapMutex.Unlock()
	return app_context.userId_websocket_map[userId]
}

type Message struct {
	From         string   `json:"from"`
	To           string   `json:"to,omitempty"`
	Body         string   `json:"body"`
	Type         string   `json:"type"`
	GroupID      string   `json:"group_id,omitempty"`
	GroupName    string   `json:"group_name,omitempty"`
	GroupMembers []string `json:"group_members,omitempty"`
}

type Group struct {
	Creater      string   `json:"creater"`
	GroupID      string   `json:"group_id,omitempty"`
	GroupName    string   `json:"group_name,omitempty"`
	GroupMembers []string `json:"group_members,omitempty"`
}

type AppContext struct {
	userId_websocket_map map[string]*websocket.Conn
	redisClient          *redis.Client
	ctx                  context.Context
	mapMutex             sync.Mutex
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
	defer app_context.RemoveConnection(userId)
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
	c := app_context.GetConnection(userId)
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
	app_context.AddConnection(userId, c)
	message_loop_handler(userId)
}

func groupHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to /create-group endpoint. Method: %s", r.Method)
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var group Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	groupJSON, err := json.Marshal(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app_context.redisClient.Set(context.Background(), group.GroupName, groupJSON, 0).Err()
	groupString, err := app_context.redisClient.Get(app_context.ctx, group.GroupName).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("group json %s\n", groupString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func main() {
	init_app_context()
	go startSubscriber(context.Background(), app_context.redisClient, "chat-message")

	//go func() {
	//	if err := logChannelMessages(context.Background(), app_context.redisClient, "chat-message"); err != nil {
	//		log.Printf("Logger error: %v", err)
	//	}
	//}()

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	http.Handle("/", webSocketHandler)
	http.HandleFunc("/create-group", groupHandler)
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

		if message.Type == "Message" {
			writeIndividualMessage(message)
		} else if message.Type == "GroupMessage" {

			// Fetch group info from Redis
			groupJson, err := rdb.Get(app_context.ctx, message.GroupName).Result()
			if err != nil {
				log.Printf("Error getting group json: %v\n", err)
				continue
			}

			var group Group
			err = json.Unmarshal([]byte(groupJson), &group)
			if err != nil {
				log.Printf("Error unmarshalling group: %v\n", err)
				continue
			}

			message.GroupMembers = group.GroupMembers

			writeGroupMessage(message)
		}

	}
}

func writeGroupMessage(message Message) {
	for _, member := range message.GroupMembers {
		c := app_context.GetConnection(member)
		if c == nil {
			log.Printf("Could not find connection to group member %s\n", member)
			continue
		}

		writeMessage := Message{
			From:      message.From,
			To:        member,
			Body:      message.Body,
			Type:      message.Type,
			GroupID:   message.GroupID,
			GroupName: message.GroupName,
		}

		err := c.WriteJSON(writeMessage)
		if err != nil {
			log.Printf("Error writing message: %v\n", err)
		}

	}
}

func writeIndividualMessage(message Message) {

	c := app_context.GetConnection(message.To)

	if c == nil {
		log.Printf("WebSocket connection not found for user: %s", message.To)

	} else {
		log.Printf("Found active WebSocket connection for user: %s", message.To)
	}

	writeMessage := Message{
		From: message.From,
		To:   message.To,
		Body: message.Body,
	}

	err := c.WriteJSON(writeMessage)
	if err != nil {
		log.Println("Error writing message:", err)
	}
}

func logChannelMessages(ctx context.Context, client *redis.Client, channelName string) error {
	// Subscribe to the channel
	pubsub := client.Subscribe(ctx, channelName)
	defer pubsub.Close()

	// Wait for confirmation of subscription
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("error subscribing to channel: %v", err)
	}

	// Get the message channel
	ch := pubsub.Channel()

	log.Printf("Starting to log messages from channel: %s\n", channelName)

	// Continuously read messages
	for msg := range ch {
		// Try to parse the message as JSON
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			// If not JSON, log as raw message
			log.Printf("Raw Message - Channel: %s, Payload: %s\n",
				msg.Channel, msg.Payload)
			continue
		}

		// Log parsed message
		log.Printf("Message - From: %s, To: %s, Body: %s\n",
			message.From,
			message.To,
			message.Body)
	}

	return nil
}
