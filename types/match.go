package types

import (
	"time"
)

type Match struct {
	UUID              *string    `json:"uuid,omitempty" bson:"uuid,omitempty"`
	UserAUUID         *string    `json:"userAUUID"`
	UserBUUID         *string    `json:"userBUUID"`
	DateCreated       *time.Time `json:"dateCreated"`
	MessageObjectUUID *string    `json:"messageObjectUUID"`
}

type MatchResponse struct {
	Token *string `json:"token"`
	Match *Match  `json:"match"`
}

type MatchRequest struct {
	UUID      *string `json:"uuid,omitempty"`
	UserAUUID *string `json:"userAUUID,omitempty"`
	UserBUUID *string `json:"userBUUID,omitempty"`
	Match     *Match  `json:"match,omitempty"`
}

type MeetingPlace struct {
	UUID        *string  `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Description *string  `json:"description"`
	Price       *string  `json:"price"`
	// other fields here
}
