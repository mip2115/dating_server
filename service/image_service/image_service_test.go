package imageservice

import (

	//"github.com/kama/server/types"

	"testing"

	"github.com/stretchr/testify/suite"
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

}
func (suite *ImageTestSuite) SetupTest() {

}

func (suite *ImageTestSuite) TearDownAllSuite() {

}

func (suite *ImageTestSuite) TearDownTest() {
}

func TestImageTestSuite(t *testing.T) {
	suite.Run(t, new(ImageTestSuite))
}
