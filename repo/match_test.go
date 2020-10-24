package repo

import (
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"

	//"github.com/kama/server/types"

	"testing"

	"code.mine/dating_server/DB"
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
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	match := &types.Match{}

	err = SaveMatch(match)
	suite.Require().NoError(err)

}

func (suite *MatchTestSuite) TestDeleteMatchByMatchUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	matchUUID := mapping.StrToPtr("match-uuid")
	match := &types.Match{
		UUID: matchUUID,
	}
	err = SaveMatch(match)
	suite.Require().NoError(err)

	err = DeleteMatchByMatchUUID(matchUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(match)

	match, err = GetMatchByMatchUUID(matchUUID)
	suite.Require().NoError(err)
	suite.Require().Nil(match)
}

func (suite *MatchTestSuite) TestGetMatchByMatchUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	match := &types.Match{
		UUID: mapping.StrToPtr("some-uuid"),
	}
	err = SaveMatch(match)
	suite.Require().NoError(err)

	m, err := GetMatchByMatchUUID(mapping.StrToPtr("some-uuid"))
	suite.Require().NoError(err)
	suite.Require().NotNil(m)

	m, err = GetMatchByMatchUUID(mapping.StrToPtr("uuid-not-found"))
	suite.Require().NoError(err)
	suite.Require().Nil(m)
}

func (suite *MatchTestSuite) TestCreateTrackedLike() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	trackedLike := &types.TrackedLike{}
	tl, err := CreateTrackedLike(trackedLike)
	suite.Require().NoError(err)
	suite.Require().NotNil(tl)
}
func (suite *MatchTestSuite) TestDeleteTrackedLike() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	matchUUID := "some-match-uuid"
	trackedLike := &types.TrackedLike{
		MatchUUID: &matchUUID,
	}
	tl, err := CreateTrackedLike(trackedLike)
	suite.Require().NoError(err)
	suite.Require().NotNil(tl)

	err = DeleteTrackedLikeByMatchUUID(trackedLike.MatchUUID)
	suite.Require().NoError(err)

	m, err := GetMatchByMatchUUID(&matchUUID)
	suite.Require().NoError(err)
	suite.Require().Nil(m)
}

func TestMatchTestSuite(t *testing.T) {

	defer DB.Disconnect()
	suite.Run(t, new(MatchTestSuite))
}
