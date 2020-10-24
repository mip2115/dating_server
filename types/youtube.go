package types

// TopicDetails -
type TopicDetails struct {
	TopicCategories []string `json: "topicCategories, omitempty"`
}

// Statistics -
type Statistics struct {
	ViewCount     string `json: "viewCount, omitempty"`
	LikeCount     string `json: "likeCount, omitempty"`
	DislikeCount  string `json: "dislikeCount, omitempty"`
	FavoriteCount string `json: "favoriteCount, omitempty"`
	CommentCount  string `json: "commentCount, omitempty"`
}

// Localized -
type Localized struct {
	Description string `json: "description, omitempty"`
}

// Snippet -
type Snippet struct {
	Title       string    `json: "title, omitempty"`
	Description string    `json: description, omitempty"`
	Tags        []string  `json: "tags, omitempty"`
	CategoryID  int       `json: "categoryId, omitempty"`
	Localized   Localized `json: "localized, omitempty"`
}

// VideoID -
type VideoID struct {
	VideoID *string `json: "videoId, omitempty"`
}

// IDItem -
type IDItem struct {
	ID *VideoID `json: "id, omitempty"`
}

// DetailsItem -
type DetailsItem struct {
	Kind         string       `json: "kind, omitempty"`
	Etag         string       `json: "etag, omitempty"`
	ID           *string      `json: "id, omitempty"`
	Snippet      Snippet      `json: "snippet, omitempty"`
	Statistics   Statistics   `json: "statistics, omitempty"`
	TopicDetails TopicDetails `json: "topicDetails, omitempty"`
}

// VideoDetailsResponse -
type UserVideoItem struct {
	UserUUID *string       `json: "userUuid"`
	UUID     *string       `json: "uuid"`
	Items    []DetailsItem `json: "items, omitempty"`
}

// VideoIDResponse -
type VideoIDResponse struct {
	Items []IDItem `json: "items, omitempty"`
}
