package youtubeService

import (
	"testing"

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

func TestYoutubeServiceTestSuite(t *testing.T) {

	suite.Run(t, new(YoutubeServiceTestSuite))
}
