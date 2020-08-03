package image_service

import (
	"context"
	"strings"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO – figure out more graceful way to deal with rank
func CreateImage(imageToUpload *types.Image) (*types.Image, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	imageToUpload.UUID = mapping.StrToPtr(uuid.String())
	// if the current rank exists then just delete that
	c, err := DB.GetCollection("images")
	if err != nil {
		return nil, err
	}

	var foundImage *types.Image
	res := c.FindOne(context.Background(), bson.M{"user_uuid": mapping.StrToV(imageToUpload.UserUUID), "rank": mapping.Int64ToV(imageToUpload.Rank)})
	if res.Err() != nil && !strings.Contains(res.Err().Error(), "no documents in result") {
		return nil, res.Err()
	}
	res.Decode(foundImage)
	if foundImage != nil {
		err = DeleteImage(foundImage.UUID)
		if err != nil {
			return nil, err
		}
	}
	c.InsertOne(context.Background(), imageToUpload)
	return imageToUpload, nil
}

func GetImagesByUserUUID(uuid *string) ([]*types.Image, error) {
	c, err := DB.GetCollection("images")
	if err != nil {
		return nil, err
	}
	results := []*types.Image{}
	res := c.FindOne(context.Background(), bson.M{"user_uuid": mapping.StrToV(uuid)})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(results)
	return results, nil
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
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": mapping.StrToV(imageUUID)})
	if err != nil {
		return err
	}
	return nil
}
