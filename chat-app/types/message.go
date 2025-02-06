package types

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

func NewMessage(message *Message) *Message {
	return &Message{
		From:      message.From,
		Body:      message.Body,
		Type:      message.Type,
		GroupID:   message.GroupID,
		GroupName: message.GroupName,
	}
}
