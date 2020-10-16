package repo

import (
	"context"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMatch -
func CreateMatch(m *types.Match) (*types.Match, error) {
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
	return m, nil
}

func DeleteMatchByMatchUUID(matchUUID *string) error {
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

func DeleteTrackedLikeByMatchUUID(matchUUID *string) error {
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"match_uuid": mapping.StrToV(matchUUID)})
	if err != nil {
		return err
	}
	return nil
}