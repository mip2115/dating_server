package repo

import (
	"context"
	"errors"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsersByFilter(filters *bson.M, options *options.FindOptions) ([]*types.User, error) {
	if filters == nil {
		return nil, errors.New("filters is nil")
	}
	if options == nil {
		return nil, errors.New("options is nil")
	}
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	cursor, err := c.Find(context.Background(), filters, options)
	users := []*types.User{}
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func GetVideosByAllUserUUIDs(userUUIDs []*string) ([]*types.UserVideoItem, error) {
	if userUUIDs == nil {
		return nil, errors.New("no userUUIDs provided")
	}

	c, err := DB.GetCollection("videos")
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"user_uuid": bson.M{
			"$in": userUUIDs,
		},
	}
	cursor, err := c.Find(context.Background(), filter)
	videos := []*types.UserVideoItem{}
	if err = cursor.All(context.Background(), &videos); err != nil {
		return nil, err
	}
	return videos, nil
}

// GetVideosByUserUUID -
func GetVideosByUserUUID(userUUID *string) ([]*types.UserVideoItem, error) {
	if userUUID == nil {
		return nil, errors.New("no userUUID provided")
	}

	c, err := DB.GetCollection("videos")
	if err != nil {
		return nil, err
	}

	videos := []*types.UserVideoItem{}
	cursor, err := c.Find(context.Background(), bson.M{"user_uuid": *userUUID})
	err = cursor.All(context.Background(), videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
