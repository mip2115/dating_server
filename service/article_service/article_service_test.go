package article_service

import (
	"io/ioutil"
	"log"
	"testing"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/aws"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type ArticleServiceTestSuite struct {
	suite.Suite
}

func (suite *ArticleServiceTestSuite) SetupSuite() {

}
func (suite *ArticleServiceTestSuite) SetupTest() {

}

func (suite *ArticleServiceTestSuite) TearDownAllSuite() {

}

func (suite *ArticleServiceTestSuite) TearDownTest() {
}

func (suite *ArticleServiceTestSuite) TestAddSynsForText() {
	content, err := ioutil.ReadFile("../../tests/articles/racism_3.txt")
	suite.NoError(err)

	contentAsString := string(content)

	err = AddSynsForText(contentAsString)
	suite.NoError(err)
}

func TestArticleServiceTestSuite(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = aws.SetAWSConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err := DB.SetupDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Client.Disconnect(*db.Ctx)
	suite.Run(t, new(ArticleServiceTestSuite))
}
