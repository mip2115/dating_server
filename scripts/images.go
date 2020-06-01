package scripts

import (
	"context"
	"github.com/kama/server/DB"
	"github.com/kama/server/mapping"
	"github.com/kama/server/types"
)

func CreateImages() error {
	err := DropCollectionImages()
	if err != nil {
		return err
	}
	image1 := types.Image{
		UUID:     mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e621"),
		UserUUID: mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fz"),
		Link:     mapping.StrToPtr("Link-to-image"),
		Rank:     mapping.IntToPtr(1),
		Key:      mapping.StrToPtr("Image-key"),
	}
	c, err := DB.GetCollection("images")
	if err != nil {
		return err
	}
	_, err = c.InsertMany(
		context.Background(),
		[]interface{}{
			image1,
		})
	if err != nil {
		return err
	}
	return nil
}

func DropCollectionImages() error {
	c, err := DB.GetCollection("images")
	if err != nil {
		return err
	}
	err = c.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
