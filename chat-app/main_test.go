package main

import (
	"github.com/gorilla/websocket"
	"go_proj/runtime"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	redis := runtime.NewRedisClient(runtime.ChatMessageChannel)
	appContext := runtime.NewAppContext(redis)
	startApplication(appContext)

	serverURL := "ws://localhost:8100/ws"
	log.Printf("Connecting to %s...", serverURL)
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("userId", "sangeet")

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	err = conn.Close()
	time.Sleep(time.Millisecond * 5)

	if err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
	appContext.Clean()
}
