package entities

import (
	"time"
)

type MessageInput struct {
	FromUser     string   `json:"from_user" db:"from_user"`
	ToUser       string   `json:"to_user,omitempty" db:"to_user"`
	Body         string   `json:"body" db:"body"`
	MessageType  string   `json:"message_type" db:"message_type"`
	GroupID      string   `json:"group_id,omitempty" db:"group_id"`
	GroupName    string   `json:"group_name,omitempty" db:"group_name"`
	GroupMembers []string `json:"group_members,omitempty" db:"group_members"`
}

type MessageResponse struct {
	ID           int64     `json:"id" db:"id"`
	FromUser     string    `json:"from_user" db:"from_user"`
	ToUser       string    `json:"to_user,omitempty" db:"to_user"`
	Body         string    `json:"body" db:"body"`
	MessageType  string    `json:"message_type" db:"message_type"`
	GroupID      string    `json:"group_id,omitempty" db:"group_id"`
	GroupName    string    `json:"group_name,omitempty" db:"group_name"`
	GroupMembers []string  `json:"group_members,omitempty" db:"group_members"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
