package userservice

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"code.mine/dating_server/mapping"
	"code.mine/dating_server/repo"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
)

type UserController struct {
	repo repo.Repo
}

func New(
	repo repo.Repo,
) *UserController {
	return &UserController{
		repo: repo,
	}
}

// make sure to also set up google/fb auth

// CreateUser -
func (c *UserController) CreateUser(userRequest *types.CreateUserRequest) (*string, error) {
	if userRequest.Email == nil {
		return nil, errors.New("need email to create user")
	}
	if userRequest.Password == nil {
		return nil, errors.New("need password to create user")
	}
	if userRequest.Password == nil || mapping.StrToV(userRequest.Password) != mapping.StrToV(userRequest.PasswordConfirm) {
		return nil, errors.New("password and password confirm do not match")
	}

	res, err := c.repo.GetUserByEmail(userRequest.Email)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return nil, errors.New("User exists")
	}

	user := &types.User{}
	user.Password = userRequest.Password
	user.Email = userRequest.Email
	user.FutureDates = []*string{}
	user.PastDates = []*string{}
	user.Matches = []*string{}
	user.UsersLikedMe = []*string{}
	user.RecentlyMatched = []*string{}
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
func (c *UserController) LikeProfile(userGettingLiked *string, userPerformingLike *string) (*types.TrackedLike, error) {
	if userGettingLiked == nil {
		return nil, errors.New("need userGettingLiked to perform like")
	}
	if userPerformingLike == nil {
		return nil, errors.New("need userPerformingLike to perform like")
	}
	// first check if there is a trackedLike
	// if there is, then just update it.
	//
	// c, err := DB.GetCollection("trackedLike")
	// if err != nil {
	// 	return nil, err
	// }
	var trackedLike *types.TrackedLike
	trackedLike, err := c.repo.GetTrackedLikeByUserUUID(userGettingLiked, userPerformingLike)
	if err != nil {
		return nil, err
	}
	if trackedLike == nil {

		newUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		trackedLike = &types.TrackedLike{
			UUID:                   mapping.StrToPtr(newUUID.String()),
			UserPerformingLikeUUID: userPerformingLike,
			UserLikedUUID:          userGettingLiked,
		}
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
		UUID:        mapping.StrToPtr(newMatchUUID.String()),
		UserOneUUID: userPerformingLike,
		UserTwoUUID: userGettingLiked,
		DateCreated: &t,
		DateUpdated: &t,
	}
	updateParams := bson.M{
		"matchUUID": newMatchUUID.String(),
	}
	filter := bson.M{
		"uuid": mapping.StrToV(trackedLike.UUID),
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

// LoginUser -
func (c *UserController) LoginUser(email, password *string) (*types.User, error) {
	if email == nil {
		return nil, errors.New("need email to create user")
	}
	if password == nil {
		return nil, errors.New("need password to create user")
	}
	checkedUser, err := c.repo.CheckUserLoginPasswordByEmail(email, password)
	if err != nil {
		return nil, err
	}
	return checkedUser, nil
}

// handle reseting ages differently because it should retrigger a search

// UpdateUser -
func (c *UserController) UpdateUser(user *types.User) error {
	fieldsToUpdate := []bson.M{}
	if user.Email != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"email": *user.Email})
	}
	if user.Mobile != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"mobile": *user.Mobile})
	}
	if user.Gender != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"gender": *user.Gender})
	}
	if user.Drink != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"drink": *user.Drink})
	}
	if user.Smoke != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"smoke": *user.Smoke})
	}
	if user.Job != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"job": *user.Job})
	}
	if user.University != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"university": *user.University})
	}
	if user.Job != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"job": *user.Job})
	}
	if user.Politics != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"politics": *user.Politics})
	}
	if user.Religion != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"religion": *user.Religion})
	}
	if user.Hometown != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"hometown": *user.Hometown})
	}
	if user.PartnerGender != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"partnerGender": *user.PartnerGender})
	}
	if user.MeetingAddress != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"meetingAddress": *user.MeetingAddress})
	}
	if user.City != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"city": *user.City})
	}
	if user.Purpose != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.M{"purpose": *user.Purpose})
	}

	err := c.repo.UpdateUserByUUID(user.UUID, fieldsToUpdate)
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

// DeleteUserByUUID -
func (c *UserController) DeleteUserByUUID(userUUID *string) error {
	err := c.repo.DeleteUserByUUID(userUUID)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByUUID -
func (c *UserController) GetUserByUUID(userUUID *string) (*types.User, error) {
	user, err := c.repo.GetUserByUUID(userUUID)
	if err != nil {
		return nil, err
	}
	return user, nil
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
