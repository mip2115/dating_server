package matchservice

import (
	"code.mine/dating_server/mapping"
	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MatchTestSuite struct {
	suite.Suite
}

func (suite *MatchTestSuite) SetupSuite() {

}
func (suite *MatchTestSuite) SetupTest() {

}

func (suite *MatchTestSuite) TearDownAllSuite() {

}

func (suite *MatchTestSuite) TearDownTest() {
}

func (suite *MatchTestSuite) TestCreateMatch() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	mockRepo.EXPECT().CreateMatch(gomock.Any()).Return(nil)
	matchController := MatchController{
		repo: mockRepo,
	}

	match := &types.Match{
		UserOneUUID: mapping.StrToPtr("user-one-uuid"),
		UserTwoUUID: mapping.StrToPtr("user-two-uuid"),
	}

	createdMatch, err := matchController.CreateMatch(match)
	suite.Require().NoError(err)
	suite.Require().NotNil(createdMatch)

	match = &types.Match{}

	createdMatch, err = matchController.CreateMatch(match)
	suite.Require().Error(err)
	suite.Require().Nil(createdMatch)
}

func (suite *MatchTestSuite) TestDeleteMatch() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)

	mockRepo.EXPECT().DeleteMatchByMatchUUID(gomock.Any()).Return(nil)
	mockRepo.EXPECT().DeleteTrackedLikeByMatchUUID(gomock.Any()).Return(nil)
	matchController := MatchController{
		repo: mockRepo,
	}

	someUUID := "some-uuid"
	err := matchController.DeleteMatch(&someUUID)
	suite.Require().NoError(err)
}
