package repo

import (
	"code.mine/dating_server/factory"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"github.com/kama/server/types"

	"testing"

	"code.mine/dating_server/DB"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite
}

func (suite *UsersTestSuite) SetupSuite() {

}
func (suite *UsersTestSuite) SetupTest() {

}

func (suite *UsersTestSuite) TearDownAllSuite() {

}

func (suite *UsersTestSuite) TearDownTest() {
}

func (suite *UsersTestSuite) TestGetTrackedLikeByUserUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	matchUUID := "match-uuid"
	userOneUUID := "user-one-uuid"
	userTwoUUID := "user-two-uuid"
	uuid := "tracked-like-uuid"
	trackedLike := &types.TrackedLike{
		UUID:                   &uuid,
		MatchUUID:              &matchUUID,
		UserPerformingLikeUUID: &userOneUUID,
		UserLikedUUID:          &userTwoUUID,
	}
	tl, err := CreateTrackedLike(trackedLike)
	suite.Require().NoError(err)
	suite.Require().NotNil(tl)
	tl, err = GetTrackedLikeByUserUUID(&userOneUUID, &userTwoUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(tl)
}

func (suite *UsersTestSuite) TestCreateUser() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	user := &types.User{
		UUID:     mapping.StrToPtr("user-uuid"),
		Password: mapping.StrToPtr("password"),
	}
	retUUID, err := CreateUser(user)
	suite.Require().NoError(err)
	suite.Require().Equal(retUUID, user.UUID)
}

func (suite *UsersTestSuite) TestUpdateTrackedLikeByUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	matchUUID := "match-uuid"
	userOneUUID := "user-one-uuid"
	userTwoUUID := "user-two-uuid"
	uuid := "tracked-like-uuid"
	trackedLike := &types.TrackedLike{
		UUID:                   &uuid,
		UserPerformingLikeUUID: &userOneUUID,
		UserLikedUUID:          &userTwoUUID,
	}
	tl, err := CreateTrackedLike(trackedLike)
	suite.Require().NoError(err)
	suite.Require().NotNil(tl)

	updateParams := bson.M{
		"matchUUID": matchUUID,
	}
	filter := bson.M{
		"uuid": uuid,
	}

	err = UpdateTrackedLikeByUUID(&uuid, filter, updateParams)
	suite.Require().NoError(err)

	trackedLike, err = GetTrackedLikeByUserUUID(&userTwoUUID, &userOneUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(trackedLike)
	suite.Require().Equal(*trackedLike.MatchUUID, matchUUID)
}

func (suite *UsersTestSuite) TestGetUsersByFilter() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("male")
		user.PartnerGender = mapping.StrToPtr("male")
		user.Zipcode = mapping.StrToPtr("10025")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("male")
		user.PartnerGender = mapping.StrToPtr("female")
		user.Zipcode = mapping.StrToPtr("10025")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("female")
		user.PartnerGender = mapping.StrToPtr("female")
		user.Zipcode = mapping.StrToPtr("10025")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("female")
		user.PartnerGender = mapping.StrToPtr("female")
		user.Zipcode = mapping.StrToPtr("10024")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("female")
		user.PartnerGender = mapping.StrToPtr("male")
		user.Zipcode = mapping.StrToPtr("10029")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	// these are the only ones that should match
	for i := 0; i < 5; i++ {
		user := factory.NewUser()
		user.Gender = mapping.StrToPtr("female")
		user.PartnerGender = mapping.StrToPtr("male")
		user.Zipcode = mapping.StrToPtr("10025")
		_, err = CreateUser(user)
		suite.Require().NoError(err)
	}

	// show me only users who have my zip code,
	// who are looking for a male
	// and who are female
	filters := &bson.M{
		"zipcode":       "10025",
		"partnerGender": "male",
		"gender":        "female",
	}
	opts := options.Find()
	opts.SetLimit(int64(10))

	users, err := GetUsersByFilter(filters, opts)
	suite.Require().NoError(err)
	suite.Require().NotNil(users)
	suite.Require().Equal(5, len(users))

}

func (suite *UsersTestSuite) TestGetUserByUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	user := factory.NewUser()
	userUUID, err := CreateUser(user)
	foundUser, err := GetUserByUUID(userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(foundUser)
}

func (suite *UsersTestSuite) TestGetUserEmail() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	user := factory.NewUser()
	_, err = CreateUser(user)
	foundUser, err := GetUserByEmail(user.Email)
	suite.Require().NoError(err)
	suite.Require().NotNil(foundUser)

	foundUser, err = GetUserByEmail(mapping.StrToPtr("different-email"))
	suite.Require().NoError(err)
	suite.Require().Nil(foundUser)
}

func (suite *UsersTestSuite) TestDeleteUserByUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	user := factory.NewUser()
	userUUID, err := CreateUser(user)
	foundUser, err := GetUserByUUID(userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(foundUser)

	err = DeleteUserByUUID(foundUser.UUID)
	suite.Require().NoError(err)

	foundUser, err = GetUserByUUID(foundUser.UUID)
	suite.Require().NoError(err)
	suite.Require().Nil(foundUser)
}

func (suite *UsersTestSuite) TestUpdateUserByUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	fieldsToUpdate := &bson.M{
		"mobile":     "8605555555",
		"first_name": "joey",
	}

	user := factory.NewUser()
	userUUID, err := CreateUser(user)
	err = UpdateUserByUUID(userUUID, fieldsToUpdate)
	suite.Require().NoError(err)

	foundUser, err := GetUserByUUID(user.UUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(foundUser)
	suite.Require().Equal("8605555555", *foundUser.Mobile)
	suite.Require().Equal("joey", *foundUser.FirstName)
}

func (suite *UsersTestSuite) TestCheckUserLoginPasswordByEmail() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	user := factory.NewUser()
	user.Password = mapping.StrToPtr("password")
	userUUID, err := CreateUser(user)
	suite.Require().NoError(err)
	suite.Require().NotNil(userUUID)

	foundUser, err := GetUserByUUID(userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(foundUser)
	u, err := CheckUserLoginPasswordByEmail(foundUser.Email, mapping.StrToPtr("password"))
	suite.Require().NoError(err)
	suite.Require().NotNil(u)

}
func TestUsersTestSuite(t *testing.T) {

	defer DB.Disconnect()
	suite.Run(t, new(UsersTestSuite))
}
