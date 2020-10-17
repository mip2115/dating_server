package userservice

import (
	"testing"

	"code.mine/dating_server/mapping"
	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
}

func (suite *UserTestSuite) SetupSuite() {

}
func (suite *UserTestSuite) SetupTest() {

}

func (suite *UserTestSuite) TearDownAllSuite() {

}

func (suite *UserTestSuite) TearDownTest() {
}

func (suite *UserTestSuite) TestGetUserByUUID() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	controller := UserController{
		repo: mockRepo,
	}

	userUUID := "some-uuid"
	user := &types.User{}
	mockRepo.EXPECT().GetUserByUUID(gomock.Any()).Return(user, nil)
	_, err := controller.GetUserByUUID(&userUUID)
	suite.Require().NoError(err)

}

func (suite *UserTestSuite) TestDeleteUserByUUID() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	controller := UserController{
		repo: mockRepo,
	}
	userUUID := "some-uuid"
	mockRepo.EXPECT().DeleteUserByUUID(gomock.Any()).Return(nil)
	err := controller.DeleteUserByUUID(&userUUID)
	suite.Require().NoError(err)

}

func (suite *UserTestSuite) TestLoginUser() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	controller := UserController{
		repo: mockRepo,
	}

	checkedUser := &types.User{}
	mockRepo.EXPECT().CheckUserLoginPasswordByEmail(gomock.Any(), gomock.Any()).Return(checkedUser, nil)

	email := "some-email"
	password := "some-password"
	user, err := controller.LoginUser(&email, &password)
	suite.Require().NoError(err)
	suite.Require().NotNil(user)
}

func (suite *UserTestSuite) TestUpdateUser() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	controller := UserController{
		repo: mockRepo,
	}
	userUpdate := &types.User{
		Mobile: mapping.StrToPtr("8604716666"),
	}
	mockRepo.EXPECT().UpdateUserByUUID(gomock.Any(), gomock.Any()).Return(nil)
	err := controller.UpdateUser(userUpdate)
	suite.Require().NoError(err)
}

// also need to add in database testing here
func (suite *UserTestSuite) TestLikeUser() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)

	newlyCreatedTrackedLike := &types.TrackedLike{
		UUID:                   mapping.StrToPtr("some-uuid"),
		UserLikedUUID:          mapping.StrToPtr("Harry"),
		UserPerformingLikeUUID: mapping.StrToPtr("James"),
	}

	trackedLikeHarryLikedJames := &types.TrackedLike{
		UUID:                   mapping.StrToPtr("some-uuid"),
		UserLikedUUID:          mapping.StrToPtr("James"),
		UserPerformingLikeUUID: mapping.StrToPtr("Harry"),
	}
	mockRepo.EXPECT().GetTrackedLikeByUserUUID(gomock.Any(), gomock.Any()).
		Return(nil, nil).
		Return(trackedLikeHarryLikedJames, nil)

	mockRepo.EXPECT().CreateTrackedLike(gomock.Any()).Return(newlyCreatedTrackedLike, nil)
	mockRepo.EXPECT().UpdateTrackedLikeByUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	mockRepo.EXPECT().SaveMatch(gomock.Any()).Return(nil)
	controller := UserController{
		repo: mockRepo,
	}

	tests := []struct {
		userGettingLiked   *string
		userPerformingLike *string
		shouldError        bool
		msg                string
	}{
		{
			userGettingLiked:   mapping.StrToPtr("Harry"),
			userPerformingLike: mapping.StrToPtr("James"),
			shouldError:        false,
			msg:                "success liking user",
		},
		{
			userGettingLiked:   mapping.StrToPtr("Harry"),
			userPerformingLike: mapping.StrToPtr("James"),
			shouldError:        false,
			msg:                "harry already liked james.  create match",
		},
	}

	for _, t := range tests {
		if !t.shouldError {
			trackedLike, err := controller.LikeProfile(t.userGettingLiked, t.userPerformingLike)
			suite.Require().NoError(err, t.msg)
			suite.Require().NotNil(trackedLike, t.msg)
		} else {
			trackedLike, err := controller.LikeProfile(t.userGettingLiked, t.userPerformingLike)
			suite.Require().Error(err, t.msg)
			suite.Require().Nil(trackedLike, t.msg)
		}
	}

}

// TODO – you should verify that if you add in a user
// it also works
func (suite *UserTestSuite) TestCreateUser() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	expectedUUID := "some-uuid"

	mockRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, nil)
	mockRepo.EXPECT().CreateUser(gomock.Any()).Return(&expectedUUID, nil)

	controller := UserController{
		repo: mockRepo,
	}

	tests := []struct {
		user        *types.CreateUserRequest
		shouldError bool
		msg         string
	}{
		{
			user: &types.CreateUserRequest{
				Email:           mapping.StrToPtr("someemail"),
				Password:        mapping.StrToPtr("somepassword"),
				PasswordConfirm: mapping.StrToPtr("somepassword"),
			},
			shouldError: false,
			msg:         "success creating user",
		},
		{
			user: &types.CreateUserRequest{
				Email:    mapping.StrToPtr("someemail"),
				Password: mapping.StrToPtr("somepassword"),
			},
			shouldError: true,
			msg:         "missing password confirm",
		},
		{
			user: &types.CreateUserRequest{
				Email:           mapping.StrToPtr("someemail"),
				Password:        mapping.StrToPtr("somepassword"),
				PasswordConfirm: mapping.StrToPtr("assword"),
			},
			shouldError: true,
			msg:         "password confirm not the same",
		},
		{
			user: &types.CreateUserRequest{
				Password:        mapping.StrToPtr("somepassword"),
				PasswordConfirm: mapping.StrToPtr("somepassword"),
			},
			shouldError: true,
			msg:         "missing email",
		},
	}

	for _, t := range tests {
		if !t.shouldError {
			uuid, err := controller.CreateUser(t.user)
			suite.Require().NoError(err, t.msg)
			suite.Require().NotNil(uuid, t.msg)
		} else {
			uuid, err := controller.CreateUser(t.user)
			suite.Require().Error(err, t.msg)
			suite.Require().Nil(uuid, t.msg)
		}
	}
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
