package userservice

import (

	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	ms "code.mine/dating_server/service/match_service"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserController strct{
	repo repo.Repo
}

// make sure to also set up google/fb auth
func (c *UserController) CreateUser(user *types.User) (*string, error) {
	user.FutureDates = []*string{}
	user.PastDates = []*string{}
	user.Matches = []*string{}
	user.UsersLikedMe = []*string{}
	user.RecentlyMatched = []*string{}

	res, err := c.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return nil, errors.New("User exists")
	}
	err = verifyEmailAndPassword(user)
	if err != nil {
		return nil, err
	}
	insertedID, err := c.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

// trackedLike
// userOne likes userTwo
// so we create a tracked like
// if userTwo likesback user one, search for a tracked like that exists
// already where user one liked userTwo.
// If yes, create a match and update Connected
func LikeProfile(userGettingLiked *string, userPerformingLike *string) (*types.TrackedLike, error) {

	// first check if there is a trackedLike
	// if there is, then just update it.
	//
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return nil, err
	}
	var trackedLike *types.TrackedLike
	trackedLike, err := c.repo.GetTrackedLikeByUserUUID(userGettingLiked,userPerformingLike)
	if err != nil {
		return nil, err
	}
	if trackedLike == nil {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		trackedLike.UUID = mapping.StrToPtr(newUUID.String())
		trackedLike.UserPerformingLikeUUID = userPerformingLike
		trackedLike.UserGettingLikedUUID = userGettingLiked
		trackedLike, err := c.repo.CreateTrackedLike(trackedLike)
		if err != nil {
			return nil, err
		}
		return trackedLike, nil
	} 
	// now create a new match
	newMatchUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	newMatch := &types.Match{
		UUID: mapping.StrToPtr(newMatchUUID.String())
		UserOneUUID: UserPerformingLikeUUID,
		UserTwoUUID: UserGettingLikedUUID,
		DateCreated: &t,
		DateUpdated: &t,
		MatchUUID: newMatchUUID,
	}
	updateParams := bson.M{
		"matchUUID": mapping.StrToV(newMatchUUID),
	}
	filter := bson.M {
		"uuid":mapping.StrToV(trackedLike.UUID),
	}
	err = c.repo.UpdateTrackedLikeByUUID(trackedLike.UUID, filter, updateParams)
	if err != nil {
		return nil, err
	}
	err = c.repo.SaveMatch(newMatch)
	if err != nil {
		return nil, err
	}
	return trackedLike, nil
}

func LoginUser(user *types.User) (*types.User, error) {
	err := verifyEmailAndPassword(user)
	if err != nil {
		return nil, err
	}
	checkedUser, err := c.repo.CheckUserLoginPasswordByEmail(user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	return checkedUser, nil
}


// handle reseting ages differently because it should retrigger a search
func UpdateUser(user *types.User) error {
	fieldsToUpdate := bson.D{}
	if user.Email != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "email", Value: *user.Email})
	}
	if user.Mobile != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "mobile", Value: *user.Mobile})
	}
	if user.Gender != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "gender", Value: *user.Gender})
	}
	if user.Drink != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "drink", Value: *user.Drink})
	}
	if user.Smoke != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "smoke", Value: *user.Smoke})
	}
	if user.Job != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "job", Value: *user.Job})
	}
	if user.University != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "university", Value: *user.University})
	}
	if user.Job != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "job", Value: *user.Job})
	}
	if user.Politics != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "politics", Value: *user.Politics})
	}
	if user.Religion != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "religion", Value: *user.Religion})
	}
	if user.Hometown != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "hometown", Value: *user.Hometown})
	}
	if user.PartnerGender != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "partnerGender", Value: *user.PartnerGender})
	}
	if user.MeetingAddress != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "meetingAddress", Value: *user.MeetingAddress})
	}
	if user.City != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "city", Value: *user.City})
	}
	if user.Purpose != nil {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "purpose", Value: *user.Purpose})
	}

	err = c.repo.UpdateUserByUUID(user.UUID, fieldToUpdate)
	if err != nil {
		return err
	}
	return nil
}

// func GetAllUsers() ([]types.User, error) {
// 	c, err := DB.GetCollection("users")
// 	if err != nil {
// 		return nil, err
// 	}
// 	cursor, err := c.Find(context.Background(), bson.D{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	users := []types.User{}
// 	if err = cursor.All(context.Background(), &users); err != nil {
// 		return nil, err
// 	}
// 	return users, nil
// }

// pick up here
func DeleteUserByUUID(userUUID *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": userUUID})
	if err != nil {
		return err
	}
	return nil
}

func GetUserByUUID(userUUID *string) (*types.User, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	user := &types.User{}
	res := c.FindOne(context.Background(), bson.M{"uuid": *userUUID})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// move this to images
func SaveUserImage(userUUID *string, imgUUID *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	update := bson.M{
		"$push": bson.M{
			"images": imgUUID,
		},
	}
	_, err = c.UpdateOne(context.Background(), bson.M{"uuid": userUUID}, update)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserImage(userUUID *string, imgUUID *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	update := bson.M{
		"$pull": bson.M{
			"images": imgUUID,
		},
	}
	_, err = c.UpdateOne(context.Background(), bson.M{"uuid": userUUID}, update)
	if err != nil {
		return err
	}
	return nil
}


func getProfilesLikedUser(userUUID *string) ([]*string, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	res := c.FindOne(context.Background(), bson.M{"uuid": *userUUID})
	user := types.User{}
	res.Decode(&user)
	return user.UsersLikedMe, nil
}

func verifyInfo(user *types.User) error {
	if user.FirstName == nil {
		return errors.New("Must provide first name")
	}
	if user.DOB == nil {
		return errors.New("Must provide DOB")
	}
	if user.Mobile == nil {
		return errors.New("Must provide Mobile")
	}
	return nil
}


func verifyEmailAndPassword(user *types.User) error {
	if user.Email == nil {
		return errors.New("must provide email")
	}
	if user.Password == nil {
		return errors.New("must provide password")
	}
	return nil
}
