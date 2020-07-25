package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"context"

	"code.mine/dating_server/server/DB"
	"code.mine/dating_server/server/mapping"
	"code.mine/dating_server/server/types"
	"github.com/nu7hatch/gouuid"
)

// TODO - add in date created
func CreateLocation(loc *types.Location) (*string, error) {
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
