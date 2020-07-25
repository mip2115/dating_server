package image_service

import (
	"context"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO – figure out way to set rank
func CreateImage(userUUID *string, link *string, key *string) (*types.Image, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	imageUUID := mapping.StrToPtr(uuid.String())
	image := types.Image{
		Link:     link,
		Rank:     nil,
		UUID:     imageUUID,
		UserUUID: userUUID,
		Key:      key,
	}
	c, err := DB.GetCollection("images")
	if err != nil {
		return nil, err
	}
	c.InsertOne(context.Background(), image)
	return &image, nil
}

func GetImageByImageUUID(uuid *string) (*types.Image, error) {
	c, err := DB.GetCollection("images")
	if err != nil {
		return nil, err
	}
	img := types.Image{}
	res := c.FindOne(context.Background(), bson.M{"uuid": *uuid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(&img)
	return &img, nil
}

func DeleteImage(imageUUID *string) error {
	c, err := DB.GetCollection("images")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": *imageUUID})
	if err != nil {
		return err
	}
	return nil
}
