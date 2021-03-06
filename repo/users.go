package repo

import (
	"context"
	"errors"

	"code.mine/dating_server/mapping"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUser -
func CreateUser(user *types.User) (*string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = mapping.StrToPtr(string(pass))
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	user.UUID = mapping.StrToPtr(u.String())
	_, err = c.InsertOne(context.Background(), user) // insert the post
	if err != nil {
		return nil, err
	}
	return user.UUID, nil
}

// GetTrackedLikeByUserUUID -
func GetTrackedLikeByUserUUID(userGettingLiked, userPerformingLike *string) (*types.TrackedLike, error) {
	trackedLike := &types.TrackedLike{}
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"userPerformingLikeUUID": mapping.StrToV(userPerformingLike),
		"userLikedUUID":          mapping.StrToV(userGettingLiked),
	}

	resp := c.FindOne(context.Background(), filter)
	err = resp.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	if err = resp.Decode(trackedLike); err != nil {
		return nil, err
	}
	return trackedLike, nil
}

// UpdateTrackedLikeByUUID -
func UpdateTrackedLikeByUUID(uuid *string, filter bson.M, updateParams bson.M) error {
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return err
	}
	params := bson.M{
		"$set": updateParams,
	}
	resp, err := c.UpdateOne(context.Background(), filter, params)
	if err != nil {
		return err
	}
	if resp.ModifiedCount == 0 {
		return errors.New("no trackedLike documents modified")
	}
	return nil

}

// GetUsersByFilter -
func GetUsersByFilter(filters *bson.M, options *options.FindOptions) ([]*types.User, error) {
	if filters == nil {
		return nil, errors.New("filters is nil")
	}
	if options == nil {
		return nil, errors.New("options is nil")
	}
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	cursor, err := c.Find(context.Background(), filters, options)
	users := []*types.User{}
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUserByUUID -
func DeleteUserByUUID(uuid *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": mapping.StrToV(uuid)})
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUUID -
func GetUserByUUID(uuid *string) (*types.User, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	resp := c.FindOne(context.Background(), bson.D{{Key: "uuid", Value: mapping.StrToV(uuid)}})
	err = resp.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	err = resp.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail -
func GetUserByEmail(email *string) (*types.User, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}

	user := &types.User{}
	resp := c.FindOne(context.Background(), bson.D{{Key: "email", Value: mapping.StrToV(email)}})
	err = resp.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, resp.Err()
	}
	err = resp.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}



// UpdateUserByUUID -
func UpdateUserByUUID(uuid *string, fieldsToUpdate *bson.M) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	update := bson.M{"$set": *fieldsToUpdate}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": mapping.StrToV(uuid)},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

// CheckUserLoginPasswordByEmail -
func CheckUserLoginPasswordByEmail(email, password *string) (*types.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), ([]byte(*password)))
	if err != nil {
		return nil, err
	}
	return user, nil
}
