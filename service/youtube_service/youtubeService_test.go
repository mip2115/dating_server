package youtubeservice

import (
	"testing"

	"code.mine/dating_server/mapping"
	"go.mongodb.org/mongo-driver/mongo/options"

	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type YoutubeServiceTestSuite struct {
	suite.Suite
}

func (suite *YoutubeServiceTestSuite) SetupSuite() {

}
func (suite *YoutubeServiceTestSuite) SetupTest() {

}

func (suite *YoutubeServiceTestSuite) TearDownAllSuite() {

}

func (suite *YoutubeServiceTestSuite) TearDownTest() {
}

func (suite *YoutubeServiceTestSuite) TestGetYoutubeVideoID() {
	query := "https://www.youtube.com/watch?v=D95qIe5pLuA"
	id, err := GetYoutubeVideoID(&query)
	suite.NoError(err)
	suite.NotNil(id)
}

func (suite *YoutubeServiceTestSuite) TestGetYoutubeVideoDetails() {
	query := "https://www.youtube.com/watch?v=D95qIe5pLuA"
	id, err := GetYoutubeVideoID(&query)
	suite.NoError(err)
	suite.NotNil(id)
	response, err := GetYoutubeVideoDetails(id)
	suite.NoError(err)
	suite.NotNil(response)
}

func (suite *YoutubeServiceTestSuite) TestGetEligibleUsers() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)

	youtube := YoutubeController{
		Repo: mockRepo,
	}

	user := &types.User{
		UUID:          mapping.StrToPtr("some-uuid"),
		PartnerGender: mapping.StrToPtr("partnerGender"),
		City:          mapping.StrToPtr("city"),
	}
	returnedUsers := []*types.User{user}

	options := options.Find()
	options.SetLimit(int64(50000))

	mockRepo.EXPECT().GetUsersByFilter(gomock.Any(), gomock.Any()).Return(returnedUsers, nil)
	users, err := youtube.GetEligibleUsers(user)
	suite.Require().NoError(err)
	suite.Require().NotNil(users)
	// suite.Require().Equal(1, len(users))
}

func TestYoutubeServiceTestSuite(t *testing.T) {

	suite.Run(t, new(YoutubeServiceTestSuite))
}
