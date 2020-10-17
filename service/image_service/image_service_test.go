package imageservice

import (

	//"github.com/kama/server/types"

	"testing"

	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
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

func (suite *ImageTestSuite) TestCreateImage() {

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	mockRepo.EXPECT().CreateImage(gomock.Any()).Return(nil)

	imageController := ImageController{
		repo: mockRepo,
	}

	img := &types.Image{}
	image, err := imageController.CreateImage(img)
	suite.Require().NoError(err)
	suite.Require().NotNil(image)
}

func (suite *ImageTestSuite) TestDeleteImage() {

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	mockRepo.EXPECT().DeleteImage(gomock.Any()).Return(nil)
	imageController := ImageController{
		repo: mockRepo,
	}

	imgUUID := "some-uuid"
	err := imageController.DeleteImage(&imgUUID)
	suite.Require().NoError(err)
}

func (suite *ImageTestSuite) TestGetImagesByUserUUID() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	imagesToReturn := []*types.Image{
		&types.Image{},
	}
	mockRepo.EXPECT().GetImagesByUserUUID(gomock.Any()).Return(imagesToReturn, nil)
	imageController := ImageController{
		repo: mockRepo,
	}
	userUUID := "user-uuid"
	images, err := imageController.GetImagesByUserUUID(&userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(images)
}

func (suite *ImageTestSuite) TestGetImageByImageUUID() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	imageToReturn := &types.Image{}

	mockRepo.EXPECT().GetImageByImageUUID(gomock.Any()).Return(imageToReturn, nil)
	imageController := ImageController{
		repo: mockRepo,
	}

	uuid := "image-uuid"
	img, err := imageController.GetImageByImageUUID(&uuid)
	suite.NoError(err)
	suite.NotNil(img)
}

func TestImageTestSuite(t *testing.T) {
	suite.Run(t, new(ImageTestSuite))
}
