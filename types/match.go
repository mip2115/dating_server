package types

import (
	"time"
)

type Match struct {
	UUID        *string    `json:"uuid,omitempty" bson:"uuid,omitempty"`
	UserOneUUID *string    `json:"userOneUUID"`
	UserTwoUUID *string    `json:"userTwoUUID"`
	DateCreated *time.Time `json:"dateCreated"`
}

type MatchResponse struct {
	Token *string `json:"token"`
	Match *Match  `json:"match"`
}

type MatchRequest struct {
	UUID        *string `json:"uuid,omitempty"`
	UserOneUUID *string `json:"userOneUUID,omitempty"`
	UserTwoUUID *string `json:"userTwoUUID,omitempty"`
	Match       *Match  `json:"match,omitempty"`
}

type MeetingPlace struct {
	UUID        *string  `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Description *string  `json:"description"`
	Price       *string  `json:"price"`
	// other fields here
}

type TrackedLike struct {
	UUID         *string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	MatchUUID    *string `json:"match_uuid" bson:"match_uuid"`
	UserOneUUID  *string `json:"userOneUUID" bson:"userOneUUID"`
	UserTwoUUID  *string `json:"userTwoUUID" bson:"userTwoUUID"`
	UserOneLiked bool    `json:"userOneLiked" bson:"userOneLiked"`
	UserTwoLiked bool    `json:"userTwoLiked" bson:"userTwoLiked"`
}
