package message_service

import (
	"code.mine/dating_server/DB"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"

	"context"
	"time"

	"code.mine/dating_server/mapping"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddMessage(
	messageRequest *types.MessageRequest,
) (*types.Message, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	msg := &types.Message{}
	msg.From = messageRequest.From
	msg.To = messageRequest.To
	msg.Content = messageRequest.Content
	msg.UUID = mapping.StrToPtr(u.String())
	t := time.Now()
	msg.DateCreated = &t
	err = addMessage(messageRequest.MatchUUID, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func GetMessages(messageRequest *types.MessageRequest, nPerPage int) ([]*types.Message, error) {
	pagesToSkip := *messageRequest.Page
	msgs, err := getMessages(pagesToSkip, nPerPage, messageRequest.MatchUUID)
	if err != nil {
		return nil, err
	}

	return msgs, nil

}

func getMessages(pagesToSkip int, nPerPage int, matchUUID *string) ([]*types.Message, error) {
	options := options.Find()
	options.SetSkip(int64(nPerPage * pagesToSkip))
	options.SetSort(bson.D{{Key: "dateCreated", Value: -1}})
	options.SetLimit(int64(nPerPage))
	c, err := DB.GetCollection("messages")
	if err != nil {
		return nil, err
	}
	// remember to set limits later on
	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	msgs := []*types.Message{}
	if err = cursor.All(context.Background(), &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func addMessage(matchUUID *string, msg *types.Message) error {
	c, err := DB.GetCollection("messages")
	if err != nil {
		return err
	}
	update := bson.M{"$push": bson.M{
		"messages": msg,
	}}
	_, err = c.UpdateOne(context.Background(), bson.M{"matchUUID": *matchUUID}, update)
	if err != nil {
		return err
	}
	return nil
}
