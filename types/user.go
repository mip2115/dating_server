package types

import (
	"time"
)

// email and/or mobile?
type CreateUserRequest struct {
	Email           *string `bson:"email" json:"email"`
	Password        *string `json:"password,omitempty" bson:"password"`
	PasswordConfirm *string `json:"password,omitempty" bson:"password"`
}

type User struct {
	UserID          *string    `bson:"_id,omitempty" json:"_id,omitempty"`
	UUID            *string    `bson:"uuid" json:"uuid"`
	Email           *string    `bson:"email" json:"email"`
	Password        *string    `json:"password,omitempty" bson:"password"`
	FirstName       *string    `bson:"first_name" json:"first_name"`
	LastName        *string    `bson:"last_name" json:"last_name"`
	Mobile          *string    `bson:"mobile" json:"mobile"`
	DOB             *string    `json:"dob" bson:"dob"`
	Gender          *string    `json:"gender" bson:"gender"`
	Age             *int64     `json:"age" bson:"age"`
	Drink           *string    `json:"drink" bson:"drink"`
	Smoke           *string    `json:"smoke" bson:"smoke"`
	Job             *string    `json:"job" bson:"job"`
	University      *string    `json:"university" bson:"university"`
	Politics        *string    `json:"politics" bson:"politics"`
	Religion        *string    `bson:"religion" json:"religion"`
	Hometown        *string    `json:"hometown" bson:"hometown"`
	PartnerGender   *string    `json:"partnerGender" bson:"partnerGender"`
	Purpose         *string    `json:"purpose" bson:"purpose"`
	Lat             *float64   `json:"lat" bson:"lat"`
	Lng             *float64   `json:"lng" bson:"lng"`
	ProfileImage    *string    `json:"profile_image" bson:"profile_image"`
	Images          []*string  `json:"images" bson:"images"`
	MeetingAddress  *string    `json:"meeting_address" bson:"meeting_address"` // meeting place within radius
	LastUpdate      *time.Time `json:"last_update" bson:"last_update"`
	PastDates       []*string  `json:"past_dates" bson:"past_dates"`
	FutureDates     []*string  `json:"future_dates" bson:"future_dates"`
	Matches         []*string  `json:"matches" bson:"matches"`
	RecentlyMatched []*string  `json:"recently_matched" bson:"recently_matched"`
	BlockedUsers    []*string  `json:"blocked_users" bson:"blocked_users"`
	UsersLikedMe    []*string  `json:"userslikedme" bson:"userslikedme"`
	MinimumAge      *int64     `json:"minimumAge" bson:"minimumAge"`
	MaximumAge      *int64     `json:"maximumAge" bson:"maximumAge"`
	City            *string    `json:"city" bson:"city"`
}

type UserResponse struct {
	Token *string `json:"token,omitempty"`
	User  *User   `json:"user,omitempty"`
	Op    *bool   `json:"op,omitempty"`
}

type UserRequest struct {
	LikedProfileID *string `json:"liked_profile_id,omitempty"`
}


