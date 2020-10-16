package repo

import (
	"context"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveMessage -
func SaveMessage(msg *types.Message) error {
	c, err := DB.GetCollection("messages")
	if err != nil {
		return err
	}
	_, err = c.InsertOne(context.Background(), msg)
	if err != nil {
		return err
	}
	return nil
}

// GetMessagesByMatchUUID -
func GetMessagesByMatchUUID(pagesToSkip int, nPerPage int, matchUUID *string) ([]*types.Message, error) {
	options := options.Find()
	options.SetSkip(int64(nPerPage * pagesToSkip))
	options.SetSort(bson.D{{Key: "dateCreated", Value: -1}})
	options.SetLimit(int64(nPerPage))
	c, err := DB.GetCollection("messages")
	if err != nil {
		return nil, err
	}
	cursor, err := c.Find(context.Background(), bson.D{}, options)
	if err != nil {
		return nil, err
	}
	msgs := []*types.Message{}
	if err = cursor.All(context.Background(), &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
