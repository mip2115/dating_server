package aws_service

import (
	//"../../DB"
	//"../../DB/aws"
	"code.mine/dating_server/server/DB"
	"code.mine/dating_server/server/DB/aws"

	//"context"
	//"../../scripts"
	"github.com/joho/godotenv"
	"code.mine/dating_server/server/scripts"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"testing"
)

var (
	ExistingUserID  = "eab85cb1-0a11-47d1-890d-93015dc1e6fz"
	fileToUploadJPG = "../../testing/images/base64.txt"
	pathEnviroment  = ""
)

type AWSTestSuite struct {
	suite.Suite
}

func (suite *AWSTestSuite) SetupSuite() {
	pathEnviroment = "../../.env"
	if err := godotenv.Load(pathEnviroment); err != nil {
		log.Fatal(err)
	}
	err := aws.SetAWSConnection()
	if err != nil {
		suite.Nil(err)
	}
	_, err = DB.SetupDB()
	suite.Nil(err)
}
func (suite *AWSTestSuite) SetupTest() {
	err := scripts.LoadDB()
	suite.Nil(err)
}

func (suite *AWSTestSuite) TearDownAllSuite() {

}

func (suite *AWSTestSuite) TearDownTest() {
}

func (suite *AWSTestSuite) TestUploadAndDeleteImageToAWS() {
	err := scripts.LoadDB()
	suite.Nil(err)

	fileBytes, err := ioutil.ReadFile(fileToUploadJPG)
	suite.Nil(err)
	_, key, err := UploadImageToS3(fileBytes, &ExistingUserID)
	suite.Nil(err)
	suite.NotNil(key)

	err = DeleteImageFromS3(key)
	suite.Nil(err)
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(AWSTestSuite))
}
