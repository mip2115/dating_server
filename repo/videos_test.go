package repo

import (
	"code.mine/dating_server/factory"
	"code.mine/dating_server/mapping"

	"testing"

	"code.mine/dating_server/types"

	"code.mine/dating_server/DB"
	"github.com/stretchr/testify/suite"
)

type VideosTestSuite struct {
	suite.Suite
}

func (suite *VideosTestSuite) SetupSuite() {

}
func (suite *VideosTestSuite) SetupTest() {

}

func (suite *VideosTestSuite) TearDownAllSuite() {

}

func (suite *VideosTestSuite) TearDownTest() {
}

func (suite *VideosTestSuite) TestGetVideosByAllUserUUIDs() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	usersToFetch := []*string{}
	userOne := mapping.StrToPtr("user-one-uuid")
	userTwo := mapping.StrToPtr("user-two-uuid")
	userThree := mapping.StrToPtr("user-three-uuid")
	for i := 0; i < 5; i++ {
		vid := factory.NewUserVideoItem()
		vid.UserUUID = userOne
		err := SaveUserVideo(vid)
		suite.Require().NoError(err)
		usersToFetch = append(usersToFetch, userOne)
	}
	for i := 0; i < 5; i++ {
		vid := factory.NewUserVideoItem()
		vid.UserUUID = userTwo
		err := SaveUserVideo(vid)
		suite.Require().NoError(err)
		usersToFetch = append(usersToFetch, userTwo)
	}
	for i := 0; i < 5; i++ {
		vid := factory.NewUserVideoItem()
		vid.UserUUID = userThree
		err := SaveUserVideo(vid)
		suite.Require().NoError(err)
	}

	vids, err := GetVideosByAllUserUUIDs(usersToFetch)
	suite.Require().NoError(err)
	suite.Require().Equal(10, len(vids))

}
func (suite *VideosTestSuite) TestGetVideosByUserUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	_, err = GetVideosByUserUUID(nil)
	suite.Require().Error(err)

	vid1 := &types.UserVideoItem{
		UUID:     mapping.StrToPtr("some-uuid"),
		UserUUID: mapping.StrToPtr("user-uuid"),
	}
	err = SaveUserVideo(vid1)
	suite.Require().NoError(err)

	vid2 := &types.UserVideoItem{
		UUID:     mapping.StrToPtr("some-other-uuid"),
		UserUUID: mapping.StrToPtr("user-uuid"),
	}
	err = SaveUserVideo(vid2)
	suite.Require().NoError(err)

	vids, err := GetVideosByUserUUID(mapping.StrToPtr("user-uuid"))
	suite.Require().NoError(err)
	suite.Require().NotNil(vids)
	suite.Require().Equal(2, len(vids))

}

func (suite *VideosTestSuite) TestSaveUserVideo() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	vid := &types.UserVideoItem{
		UUID: mapping.StrToPtr("some-uuid"),
	}
	err = SaveUserVideo(vid)
	suite.Require().NoError(err)

}
func TestVideosTestSuite(t *testing.T) {

	defer DB.Disconnect()
	suite.Run(t, new(VideosTestSuite))
}
