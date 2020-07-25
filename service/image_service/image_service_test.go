package image_service

import (
	"github.com/joho/godotenv"
	"code.mine/dating_server/server/DB"
	"code.mine/dating_server/server/scripts"
	//"github.com/kama/server/types"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

var (
	ExisingImageUUID = "eab85cb1-0a11-47d1-890d-93015dc1e621"
	ExistingUserUUID = "eab85cb1-0a11-47d1-890d-93015dc1e6fz"
	fileToUploadJPG  = "../../testing/images/base64.txt"
	pathEnviroment   = ""
)

type ImageTestSuite struct {
	suite.Suite
}

func (suite *ImageTestSuite) SetupSuite() {
	pathEnviroment = "../../.env"
	if err := godotenv.Load(pathEnviroment); err != nil {
		log.Fatal(err)
	}
	_, err := DB.SetupDB()
	suite.Nil(err)
}
func (suite *ImageTestSuite) SetupTest() {
	err := scripts.LoadDB()
	suite.Nil(err)
}

func (suite *ImageTestSuite) TearDownAllSuite() {

}

func (suite *ImageTestSuite) TearDownTest() {
}

func (suite *ImageTestSuite) TestReadImageByImageUUID() {
	err := scripts.LoadDB()
	suite.Nil(err)

	img, err := GetImageByImageUUID(&ExisingImageUUID)
	suite.Nil(err)
	suite.NotNil(img)
	suite.Equal(ExisingImageUUID, *img.UUID)
}

func (suite *ImageTestSuite) TestCreateImage() {
	err := scripts.LoadDB()
	suite.Nil(err)

	userUUID := "user-uuid"
	link := "link"
	key := "key"
	img, err := CreateImage(&userUUID, &link, &key)
	suite.Nil(err)
	suite.NotNil(img)
	returnedImage, err := GetImageByImageUUID(img.UUID)
	suite.Nil(err)
	suite.NotNil(returnedImage)

}

func TestImageTestSuite(t *testing.T) {
	suite.Run(t, new(ImageTestSuite))
}
