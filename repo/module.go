package repo

import (
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo interface {
	GetUsersByFilter(filters *bson.M, options *options.FindOptions) ([]*types.User, error)
	GetVideosByAllUserUUIDs(userUUIDs []*string) ([]*types.UserVideoItem, error)
	GetVideosByUserUUID(userUUID *string) ([]*types.UserVideoItem, error)
}
