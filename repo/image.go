package repo

import (
	"context"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateImage -
func CreateImage(img *types.Image) error {
	c, err := DB.GetCollection("images")
	if err != nil {
		return err
	}
	_, err = c.InsertOne(context.Background(), img)
	if err != nil {
		return err
	}
	return nil
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
	img := &types.Image{}
	res := c.FindOne(context.Background(), bson.M{"uuid": *uuid})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(img)
	return img, nil
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
