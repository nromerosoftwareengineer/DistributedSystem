package converter

import (
	msg "go_proj/database/entities"
	"go_proj/types"
)

func NewMessageInput(message types.Message) msg.MessageInput {
	return msg.MessageInput{
		FromUser:     message.From,
		ToUser:       message.To,
		Body:         message.Body,
		MessageType:  message.Type,
		GroupID:      message.GroupID,
		GroupName:    message.GroupName,
		GroupMembers: message.GroupMembers,
	}
}
