package DB

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func (suite *DatabaseTestSuite) SetupSuite() {

}
func (suite *DatabaseTestSuite) SetupTest() {

}

func (suite *DatabaseTestSuite) TearDownAllSuite() {

}

func (suite *DatabaseTestSuite) TearDownTest() {
}

func (suite *DatabaseTestSuite) TestGetTestDB() {
	_, err := GetTestDB()
	suite.Require().NoError(err)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
