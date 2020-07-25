package types

type Image struct {
	UUID     *string `json:"uuid" bson:"uuid"`
	UserUUID *string `json:"user_uuid" bson:"user_uuid"`
	Rank     *int64  `json:"rank" bson:"rank"`
	Link     *string `json:"link" bson:"link"`
	Key      *string `json:"key" bson:"key"`
}

type ImageUploadRequest struct {
	FileBytes []byte `json:"fileBytes" bson:"fileBytes"`
	Rank      *int64 `json:"rank" bson:"rank"`
}

type ImageDeleteRequest struct {
	UUID *string `json:"uuid" bson:"uuid"`
	Key  *string `json:"key" bson:"key"`
	Rank *int64  `json:"rank" bson:"rank"`
}

type UploadImageResponse struct {
	Token             *string `json:"token,omitempty"`
	UploadedImageLink *string `json:"uploadedImageLink"`
}

type ImageDeleteResponse struct {
	Token *string `json:"token,omitempty"`
}
