package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

type Redis struct {
	client *redis.Client
	pubsub *redis.PubSub
	ctx    context.Context
}

func NewRedisClient(channelName string) *Redis {
	ctx := context.Background()
	client := GetClient(ctx)
	pubsub := GetPubSubFor(ctx, client, channelName)
	return &Redis{
		client,
		pubsub,
		ctx,
	}
}

func GetClient(ctx context.Context) *redis.Client {
	runOnContainer := os.Getenv("IS_CONTAINER_RUN")
	var address string
	if runOnContainer == "true" {
		address = "redis:6379"
		log.Println("Redis running on " + address)
	} else {
		log.Println("Redis running on " + address)
		address = "0.0.0.0:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Redis connected:", pong)
	return rdb
}

func GetPubSubFor(ctx context.Context, client *redis.Client, channelName string) *redis.PubSub {
	pubsub := client.Subscribe(ctx, channelName)
	fmt.Printf("Subscribed to channel: %s\n", channelName)
	return pubsub
}

func (r *Redis) Close() {
	err := r.pubsub.Close()
	if err != nil {
		log.Println("Error closing redis PubSub:", err)
	}
	err = r.client.Close()
	if err != nil {
		log.Println("Error closing redis client:", err)
	}
}

func (r *Redis) Publish(data []byte, channelName string) {
	err := r.client.Publish(r.ctx, channelName, data).Err()
	if err != nil {
		log.Printf("Error publishing message: %v\n", err)
	}
}

func (r *Redis) HandlePublishedMessages(mh *MessageHandler) {
	ch := r.pubsub.Channel()
	for msg := range ch {
		var message Message
		err := json.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			log.Printf("Invalid JSON message: %s\n", msg.Payload)
			continue
		}

		toUsers := []string{}
		if message.Type == "Message" {
			toUsers = append(toUsers, message.To)
		} else if message.Type == "GroupMessage" {
			toUsers = r.GetUsersInGroup(message)
		}

		if toUsers == nil {
			continue
		}
		mh.SendMessageToUsers(NewMessage(&message), toUsers)
	}
}

func (r *Redis) GetUsersInGroup(message Message) []string {
	groupJson, err := r.Get(message.GroupName)
	if err != nil {
		log.Printf("Error getting group json: %v\n", err)
		return nil
	}

	var group Group
	err = json.Unmarshal([]byte(groupJson), &group)
	if err != nil {
		log.Printf("Error unmarshalling group: %v\n", err)
		return nil
	}
	return group.GroupMembers
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *Redis) Set(key string, data []byte) {
	r.client.Set(r.ctx, key, data, 0)
}
