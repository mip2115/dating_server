package user_service

import (
	//	"../auth"
	//"../mapping"
	//"../types"

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

func CreateUser(user *types.User) (*string, error) {
	user.FutureDates = []*string{}
	user.PastDates = []*string{}
	user.Matches = []*string{}
	user.UsersLikedMe = []*string{}
	user.RecentlyMatched = []*string{}

	res, err := checkUserExists(user)
	if err != nil {
		return nil, err
	}
	if res == false {
		return nil, errors.New("User exists")
	}
	// some error checking
	err = verifyInfo(user)
	if err != nil {
		return nil, err
	}
	insertedID, err := createUser(user)
	if err != nil {
		return nil, err
	}
	return insertedID, nil
}

// TODO –
// best way to do this is just to create a "tracked_like" kind of situation
// upload a record and then you can also tell if
// they both liked each other
// always set UserOne as the first to make the like
func LikeProfile(userUUID *string, likedProfileUUID *string) (*types.TrackedLike, error) {

	// first check if there is a trackedLike
	// if there is, then just update it.
	//
	c, err := DB.GetCollection("trackedLike")
	if err != nil {
		return nil, err
	}

	trackedLike := &types.TrackedLike{}
	query := []bson.D{bson.D{{Key: "userOneUUID", Value: *userUUID}}, {{Key: "userOneUUID", Value: *likedProfileUUID}}}
	res := c.FindOne(context.Background(), bson.D{{Key: "$or", Value: query}})
	err = res.Err()

	var createTrackedLike bool
	if strings.Contains(err.Error(), "no documents in result") {
		createTrackedLike = true
	}
	if err != nil && !createTrackedLike {
		return nil, err
	}
	if createTrackedLike {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		trackedLike.UUID = mapping.StrToPtr(newUUID.String())
		trackedLike.UserOneUUID = userUUID
		trackedLike.UserTwoUUID = likedProfileUUID
		trackedLike.UserOneLiked = true
		_, err = c.InsertOne(context.Background(), trackedLike)
		if err != nil {
			return nil, err
		}
		return trackedLike, nil
	}
	res.Decode(trackedLike)
	if trackedLike.UserOneUUID == nil {
		return nil, fmt.Errorf("user one cannot be nil in tracked like %s ", mapping.StrToV(trackedLike.UUID))
	}
	if trackedLike.UserTwoUUID == nil {
		return nil, fmt.Errorf("user two should not be nil in tracked like %s ", mapping.StrToV(trackedLike.UUID))
	}
	if mapping.StrToV(trackedLike.UserOneUUID) != mapping.StrToV(likedProfileUUID) {
		return nil, fmt.Errorf("user one should be the profile we're liking in tracked like %s ", mapping.StrToV(trackedLike.UUID))
	}
	if !trackedLike.UserOneLiked {
		return nil, fmt.Errorf("user one liked should be true in tracked like %s ", mapping.StrToV(trackedLike.UUID))
	}
	if trackedLike.MatchUUID != nil {
		return nil, fmt.Errorf("match uuid in tracked like %s should be nil", mapping.StrToV(trackedLike.UUID))
	}

	fieldsToUpdate := bson.D{}
	fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "userTwoLiked", Value: true})

	// now generate the match
	match := &types.Match{}
	match.UserOneUUID = trackedLike.UserOneUUID
	match.UserTwoUUID = trackedLike.UserTwoUUID
	insertedID, err := ms.CreateMatch(match)
	if err != nil {
		return nil, err
	}
	fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "match_uuid", Value: mapping.StrToV(insertedID)})

	trackedLike.UserTwoLiked = true
	trackedLike.MatchUUID = insertedID
	update := bson.D{{Key: "$set",
		Value: fieldsToUpdate,
	}}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": mapping.StrToV(trackedLike.UUID)},
		update,
	)
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
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	res := c.FindOne(context.Background(), bson.D{{Key: "email", Value: user.Email}})
	if res.Err() != nil {
		return nil, res.Err()
	}
	returnedUser := &types.User{}
	err = bcrypt.CompareHashAndPassword([]byte(*returnedUser.Password), ([]byte(*user.Password)))
	if err != nil {
		return nil, err
	}
	return returnedUser, nil
}

// handle reseting ages differently because it should retrigger a search
func UpdateUser(user *types.User) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
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

	update := bson.D{{Key: "$set",
		Value: fieldsToUpdate,
	}}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": *user.UUID},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]types.User, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	users := []types.User{}
	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(userUUID *string) error {
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

func GetUser(userUUID *string) (*types.User, error) {
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

// remove userIDA from list of users that likes userIDB
func pullProfileFromUserLikes(userAUUID *string, userBUUID *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}

	update := bson.M{"$pull": bson.M{"userslikedme": userAUUID}}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": *userBUUID},
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func addProfileToUserLikes(profileAUUID *string, profileBUUID *string) error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	update := bson.M{"$push": bson.M{"userslikedme": *profileAUUID}}
	_, err = c.UpdateOne(
		context.Background(),
		bson.M{"uuid": *profileBUUID},
		update,
	)
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

// use and return UUID's
func createUser(user *types.User) (*string, error) {
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

func checkUserExists(user *types.User) (bool, error) {
	err := verifyEmailAndPassword(user)
	if err != nil {
		return false, err
	}
	c, err := DB.GetCollection("users")
	if err != nil {
		return false, err
	}
	resp := c.FindOne(context.Background(), bson.D{{Key: "email", Value: user.Email}})
	if resp.Err() != nil {
		return true, nil
	}
	return false, nil
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

func contains(arr []*string, tar string) bool {
	for _, v := range arr {
		if *v == tar {
			return true
		}
	}
	return false
}
