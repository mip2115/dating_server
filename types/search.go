package types

type SearchRequest struct {
	SkipValue *int64  `json:"skipValue" bson:"skipValue"`
	UserUUID  *string `json:"userUUID" bson:"userUUID"`
}
