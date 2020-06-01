package types

import (
	"time"
)

type Message struct {
	UUID        *string    `json:"uuid,omitempty" bson:"uuid,omitempty"`
	From        *string    `json:"userFrom"`
	To          *string    `json:"userTo"`
	DateCreated *time.Time `json:"dateCreated"`
	Content     *string    `json:"content"`
}

type MessagesObject struct {
	UUID      *string    `json:"uuid,omitempty" bson:"uuid,omitempty"`
	MatchUUID *string    `json:"matchUUID" bson:"matchUUID"`
	Messages  []*Message `json:"messages"`
}

type MessageRequest struct {
	From      *string `json:"userFrom"`
	To        *string `json:"userTo"`
	Content   *string `json:"content"`
	MatchUUID *string `json:"matchUUID"`
	UUID      *string `json:"uuid"`
	Page      *int    `json:"page"`
}

type MessageResponse struct {
	Message *Message `json:"message"`
	Token   *string  `json:"token"`
}

type MessagesObjectResponse struct {
	Messages []*Message `json:"messages"`
	Token    *string    `json:"token"`
}
