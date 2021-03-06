package types

import (
	"time"
)

// Article is a record of an article uploaded by the user for matching
type Article struct {
	UserUUID    *string      `json:"user_uuid" bson:"user_uuid"`
	UUID        *string      `json:"uuid" bson:"uuid"`
	Text        *string      `json:"text" bson:"text"`
	DateCreated *time.Time   `json:"date_created" bson:"date_created"`
	WordCount   *int64       `json:"word_count" bson:"word_count"`
	Summary     *TextSummary `json:"textSummary" bson:"textSummary"`
}

type TextSummary struct {
	TopRatedWords []TopRatedWord
	Keyphrases    []string
	Entities      []string
	Key           int
}

type TopRatedWord struct {
	Word  string
	Score float32
}

type WordInformation struct {
	WordList  []string  `json:"syn" bson:"syn"`
	UUID      string    `json:"uuid" bson:"uuid"`
	Word      string    `json:"word" bson:"word"`
	WordStem  string    `json:"wordStem" bson:"wordStem"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// Noun WordList `json:"noun" bson:"noun"`
}

type WordList struct {
	Word []string `json:"syn" bson:"syn"`
}
