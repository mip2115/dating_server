package awsservice

import (
	"io/ioutil"
	"log"
	"testing"

	mockGateway "code.mine/dating_server/mocks/gateway"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AWSTestSuite struct {
	suite.Suite
}

func (suite *AWSTestSuite) SetupSuite() {

}
func (suite *AWSTestSuite) SetupTest() {

}

func (suite *AWSTestSuite) TearDownAllSuite() {

}

func (suite *AWSTestSuite) TearDownTest() {
}

func (suite *AWSTestSuite) TestUploadImageToS3() {
	content, err := ioutil.ReadFile("base64.txt")
	if err != nil {
		log.Fatal(err)
	}
	userUUID := "some-user-uuid"

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockGateway := mockGateway.NewMockGateway(mockCtrl)
	mockGateway.EXPECT().UploadUserImageToS3(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	awsController := AWSController{
		gateway: mockGateway,
	}

	key, err := awsController.UploadImageToS3(content, &userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(key)
}

func (suite *AWSTestSuite) TestDeleteImageFromS3() {
	keyOfImage := "someKey"

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockGateway := mockGateway.NewMockGateway(mockCtrl)
	mockGateway.EXPECT().DeleteUserImageFromS3(gomock.Any()).Return(nil)
	awsController := AWSController{
		gateway: mockGateway,
	}

	err := awsController.DeleteImageFromS3(&keyOfImage)
	suite.Require().NoError(err)
}

func TestSearchTestSuite(t *testing.T) {
	suite.Run(t, new(AWSTestSuite))
}
