package search_service

import (
	"code.mine/dating_server/server/DB"

	//"context"
	"github.com/joho/godotenv"
	"code.mine/dating_server/server/scripts"
	"github.com/stretchr/testify/suite"
	//	"log"

	"testing"
)

var (
	ExistingUserID = "eab85cb1-0a11-47d1-890d-93015dc1e6fz"
	pathEnviroment = ""
)

type SearchTestSuite struct {
	suite.Suite
}

func (suite *SearchTestSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
	suite.Nil(err)
	/*
		pathEnviroment = "../.env"
		if err := godotenv.Load(pathEnviroment); err != nil {
			log.Fatal(err)
		}
	*/

	_, err = DB.SetupDB()
	if err != nil {
		panic(err.Error())
	}
}
func (suite *SearchTestSuite) SetupTest() {
	err := scripts.LoadDB()
	suite.Nil(err)

}

func (suite *SearchTestSuite) TearDownAllSuite() {

}

func (suite *SearchTestSuite) TearDownTest() {

}

func (suite *SearchTestSuite) TestCalculateTopMatches() {
	err := scripts.LoadDB()
	suite.Nil(err)

	users, err := CalculateTopUsers(&ExistingUserID, 0)
	suite.Nil(err)
	suite.NotNil(users)
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

// add funcitonality for scripts to load up DB
// see if you can get somethign similar to MakeTest working
// AWS?
// factory
