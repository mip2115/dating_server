package repo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
)

// GetVideosByAllUserUUIDs -
func GetVideosByAllUserUUIDs(userUUIDs []*string) ([]*types.UserVideoItem, error) {
	if userUUIDs == nil {
		return nil, errors.New("no userUUIDs provided")
	}

	c, err := DB.GetCollection("videos")
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"useruuid": bson.M{
			"$in": userUUIDs,
		},
	}
	cursor, err := c.Find(context.Background(), filter)
	videos := []*types.UserVideoItem{}
	if err = cursor.All(context.Background(), &videos); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no videos found")
		}
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
	cursor, err := c.Find(context.Background(), bson.M{"useruuid": *userUUID})
	err = cursor.All(context.Background(), &videos)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no videos found")
		}
		return nil, err
	}
	return videos, nil
}

// SaveUserVideo -
func SaveUserVideo(video *types.UserVideoItem) error {
	c, err := DB.GetCollection("videos")
	if err != nil {
		return err
	}

	_, err = c.InsertOne(context.Background(), video)
	if err != nil {
		return err
	}
	return nil
}
