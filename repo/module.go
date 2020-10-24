package repo

//go:generate mockgen -destination=mocks/repo/mock_repo.go -package=mocks code.mine/dating_server/repo Repo

import (
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// drop the entire database before you do any repo testing down here

// Repo -
type Repo interface {
	CheckUserLoginPasswordByEmail(email, password *string) (*types.User, error)
	CreateImage(img *types.Image) error
	CreateMatch(m *types.Match) (*types.Match, error)
	CreateTrackedLike(trackedLike *types.TrackedLike) (*types.TrackedLike, error)
	CreateUser(user *types.User) (*string, error)
	DeleteImage(imageUUID *string) error
	DeleteMatchByMatchUUID(matchUUID *string) error
	DeleteTrackedLikeByMatchUUID(matchUUID *string) error
	DeleteUserByUUID(uuid *string) error
	GetImageByImageUUID(uuid *string) (*types.Image, error)
	GetImagesByUserUUID(uuid *string) ([]*types.Image, error)
	GetMessagesByMatchUUID(pagesToSkip int, nPerPage int, matchUUID *string) ([]*types.Message, error)
	GetTrackedLikeByUserUUID(userGettingLiked, userPerformingLike *string) (*types.TrackedLike, error)
	GetUserByEmail(email *string) (*types.User, error)
	GetUserByUUID(uuid *string) (*types.User, error)
	GetUsersByFilter(filters *bson.M, options *options.FindOptions) ([]*types.User, error)
	GetVideosByAllUserUUIDs(userUUIDs []*string) ([]*types.UserVideoItem, error)
	GetVideosByUserUUID(userUUID *string) ([]*types.UserVideoItem, error)
	SaveMatch(newMatch *types.Match) error
	SaveMessage(msg *types.Message) error
	UpdateTrackedLikeByUUID(uuid *string, filter bson.M, updateParams bson.M) error
	UpdateUserByUUID(uuid *string, fieldsToUpdate bson.M) error
}
