package main

import (
	"context"
	"fmt"
	db "go_proj/database"
	"go_proj/database/entities"
	"log"
)

func main() {

	dbService, err := db.NewDBService()

	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer dbService.Close()

	// Create a test message
	msg := entities.MessageInput{
		FromUser:     "user3",
		ToUser:       "user2",
		Body:         "Hello!",
		MessageType:  "direct",
		GroupMembers: []string{"user1", "user2", "user3"},
	}
	// Insert the message
	messageID, err := dbService.InsertMessage(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to insert message: %v", err)
	}

	fmt.Printf("Successfully inserted message with ID: %d\n", messageID)
}
