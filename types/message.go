package types

import (
	"time"
)

type Message struct {
	UUID        *string    `json:"uuid,omitempty" bson:"uuid,omitempty"`
	From        *string    `json:"userFrom"`
	To          *string    `json:"userTo"`
	DateCreated *time.Time `json:"dateCreated"`
	DateUpdated *time.Time `json:"dateUpdated"`
	Content     *string    `json:"content"`
	MatchUUID   *string    `json:"match_uuid" bson:"match_uuid"`
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
