package types

import (
	"time"
)

// Article is a record of an article uploaded by the user for matching
type Article struct {
	UserUUID    *string    `json:"user_uuid" bson:"user_uuid"`
	UUID        *string    `json:"uuid" bson:"uuid"`
	Text        *string    `json:"text" bson:"text"`
	DateCreated *time.Time `json:"date_created" bson:"date_created"`
	WordCount   *int64     `json:"word_count" bson:"word_count"`
}

type TextSummary struct {
	TopRatedWords []TopRatedWord
	Keyphrases    []string
	Entities      []string
}

type TopRatedWord struct {
	Word  string
	Score float32
}

type WordInformation struct {
	Noun WordList `json:"noun" bson:"noun"`
}

type WordList struct {
	Word []string `json:"syn" bson:"syn"`
}
