package main

import (
	"go_proj/api"
	"go_proj/runtime"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func startApplication(appContext *runtime.AppContext) {
	go appContext.MH.HandlePublishedMessages()
	webSocketHandler := api.NewWebSocketHandler(appContext)
	httpHandler := api.NewHttpHandler(appContext, appContext.MH.RedisClient)

	http.Handle("/", webSocketHandler)
	http.HandleFunc("/create-group", httpHandler.ServeHTTP)
	go startServer()
}

func main() {
	redis := runtime.NewRedisClient(runtime.ChatMessageChannel)
	appContext := runtime.NewAppContext(redis)
	defer appContext.Clean()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	startApplication(appContext)
	<-interrupt

	log.Println("\nSignal received, shutting down...")
}

func startServer() {
	log.Print("Starting server...")
	err := http.ListenAndServe("0.0.0.0:8100", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

//func logChannelMessages(ctx context.Context, client *redis.Client, channelName string) error {
//	// Subscribe to the channel
//	pubsub := client.Subscribe(ctx, channelName)
//	defer pubsub.Close()
//
//	// Wait for confirmation of subscription
//	_, err := pubsub.Receive(ctx)
//	if err != nil {
//		return fmt.Errorf("error subscribing to channel: %v", err)
//	}
//
//	// Get the message channel
//	ch := pubsub.Channel()
//
//	log.Printf("Starting to log messages from channel: %s\n", channelName)
//
//	// Continuously read messages
//	for msg := range ch {
//		// Try to parse the message as JSON
//		var message Message
//		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
//			// If not JSON, log as raw message
//			log.Printf("Raw Message - Channel: %s, Payload: %s\n",
//				msg.Channel, msg.Payload)
//			continue
//		}
//
//		// Log parsed message
//		log.Printf("Message - From: %s, To: %s, Body: %s\n",
//			message.From,
//			message.To,
//			message.Body)
//	}
//
//	return nil
//}
