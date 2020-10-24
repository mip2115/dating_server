package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateMatch -
// func CreateMatch(m *types.Match) (*types.Match, error) {
// 	c, err := DB.GetCollection("matches")
// 	uuid, err := uuid.NewV4()
// 	if err != nil {
// 		return nil, err
// 	}
// 	m.UUID = mapping.StrToPtr(uuid.String())
// 	_, err = c.InsertOne(context.Background(), m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return m, nil
// }

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

// GetMatchByMatchUUID -
func GetMatchByMatchUUID(matchUUID *string) (*types.Match, error) {
	c, err := DB.GetCollection("matches")
	if err != nil {
		return nil, err
	}
	resp := c.FindOne(context.Background(), bson.M{"uuid": *matchUUID})
	if resp.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	err = resp.Err()
	if err != nil {
		return nil, err
	}

	match := &types.Match{}
	err = resp.Decode(match)
	if err != nil {
		return nil, err
	}
	return match, nil
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

// SaveMatch -
func SaveMatch(newMatch *types.Match) error {
	c, err := DB.GetCollection("matches")
	if err != nil {
		return err
	}

	_, err = c.InsertOne(context.Background(), newMatch)
	if err != nil {
		return err
	}
	return nil
}

// CreateTrackedLike -
func CreateTrackedLike(trackedLike *types.TrackedLike) (*types.TrackedLike, error) {
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return nil, err
	}

	_, err = c.InsertOne(context.Background(), trackedLike)
	if err != nil {
		return nil, err
	}
	return trackedLike, nil
}
