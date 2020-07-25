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
	insertedMatchUUID, err := createMatch(m)
	if err != nil {
		return nil, err
	}
	return insertedMatchUUID, nil
}

func DeleteMatch(m *types.MatchRequest) error {
	err := deleteMatch(m.UUID)
	if err != nil {
		return err
	}
	err = deleteTrackedLikeByMatchUUID(m.UUID)
	if err != nil {
		return err
	}
	return nil
}

func deleteTrackedLikeByMatchUUID(matchUUID *string) error {
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
