package match_service

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	"context"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
)

// TODO - error check to make sure youre not matching people who recently matched
// or who should be blocked
func CreateMatch(m *types.Match) (*string, error) {
	t := time.Now()
	m.DateCreated = &t
	// first create messages document
	msgObjUUID, err := createMessagesObj()
	if err != nil {
		return nil, err
	}
	m.MessageObjectUUID = msgObjUUID
	insertedMatchUUID, err := createMatch(m)
	if err != nil {
		return nil, err
	}
	err = addMatchToUser(insertedMatchUUID, m.UserAUUID)
	if err != nil {
		return nil, err
	}
	err = addMatchToUser(insertedMatchUUID, m.UserBUUID)
	if err != nil {
		return nil, err
	}
	return insertedMatchUUID, nil
}

func DeleteMatch(m *types.MatchRequest) error {
	err := removeMatchMsgs(m.Match.MessageObjectUUID)
	if err != nil {
		return err
	}
	err = deleteMatch(m.UUID)
	if err != nil {
		return err
	}
	err = removeMatchFromUser(m.UUID, m.UserBUUID)
	if err != nil {
		return err
	}
	err = removeMatchFromUser(m.UUID, m.UserAUUID)
	if err != nil {
		return err
	}
	return nil
}

func removeMatchMsgs(msgUUID *string) error {
	c, err := DB.GetCollection("messages")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": *msgUUID})
	if err != nil {
		return err
	}
	return nil
}

func createMessagesObj() (*string, error) {
	msgObj := &types.MessagesObject{}
	msgs := []*types.Message{}
	msgObj.Messages = msgs
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	msgObj.UUID = mapping.StrToPtr(uuid.String())
	c, err := DB.GetCollection("messages")
	if err != nil {
		return nil, err
	}
	_, err = c.InsertOne(context.Background(), msgObj)
	if err != nil {
		return nil, err
	}
	msgObjUUID := uuid.String()
	return &msgObjUUID, nil
}

func deleteMatch(matchUUID *string) error {
	c, err := DB.GetCollection("matches")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": *matchUUID})
	if err != nil {
		return err
	}
	return nil
}

func createMatch(m *types.Match) (*string, error) {
	c, err := DB.GetCollection("matches")
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	m.UUID = mapping.StrToPtr(uuid.String())
	_, err = c.InsertOne(context.Background(), m)
	if err != nil {
		return nil, err
	}
	return m.UUID, nil
}

func addMatchToUser(matchUUID *string, userUUID *string) error {
	update := bson.M{"$push": bson.M{"matches": *matchUUID}}
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": *userUUID},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func removeMatchFromUser(matchUUID *string, userUUID *string) error {
	update := bson.M{"$pull": bson.M{"matches": *matchUUID}}
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": *userUUID},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}
