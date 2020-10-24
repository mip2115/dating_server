package repo

import (
	"code.mine/dating_server/factory"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"

	//"github.com/kama/server/types"

	"testing"

	"code.mine/dating_server/DB"
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

func (suite *ImageTestSuite) TestGetImagesByUserUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	images := []*types.Image{}
	userUUID := mapping.StrToPtr("user-uuid")
	for i := 0; i < 5; i++ {
		img := factory.NewImage()
		img.UserUUID = userUUID
		images = append(images, img)

		err := CreateImage(img)
		suite.Require().NoError(err)
	}

	imagesOfUser, err := GetImagesByUserUUID(userUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(imagesOfUser)
	suite.Require().Equal(len(images), 5)

}
func (suite *ImageTestSuite) TestCreateImage() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	image := &types.Image{
		UUID:     mapping.StrToPtr("some-uuid"),
		UserUUID: mapping.StrToPtr("user-uuid"),
		Link:     mapping.StrToPtr("some-link"),
		Key:      mapping.StrToPtr("some-key"),
	}
	err = CreateImage(image)
	suite.Require().NoError(err)

	queriedImage, err := GetImageByImageUUID(mapping.StrToPtr("some-uuid"))
	suite.Require().NoError(err)
	suite.Require().Equal(queriedImage.UUID, image.UUID)
}

func (suite *ImageTestSuite) TestGetImageByImageUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	image := &types.Image{
		UUID:     mapping.StrToPtr("some-uuid"),
		UserUUID: mapping.StrToPtr("user-uuid"),
		Link:     mapping.StrToPtr("some-link"),
		Key:      mapping.StrToPtr("some-key"),
	}
	err = CreateImage(image)
	suite.Require().NoError(err)

	queriedImage, err := GetImageByImageUUID(mapping.StrToPtr("some-uuid"))
	suite.Require().NoError(err)
	suite.Require().Equal(queriedImage.UUID, image.UUID)

	queriedImage, err = GetImageByImageUUID(mapping.StrToPtr("some-other-uuid"))
	suite.Require().Nil(queriedImage)
	suite.Require().Nil(err)
}

func (suite *ImageTestSuite) TestDeleteImage() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	imageUUID := mapping.StrToPtr("some-uuid")
	image := &types.Image{
		UUID:     imageUUID,
		UserUUID: mapping.StrToPtr("user-uuid"),
		Link:     mapping.StrToPtr("some-link"),
		Key:      mapping.StrToPtr("some-key"),
	}
	err = CreateImage(image)
	suite.Require().NoError(err)

	queriedImage, err := GetImageByImageUUID(imageUUID)
	suite.Require().NoError(err)
	suite.Require().NotNil(queriedImage)

	err = DeleteImage(imageUUID)
	suite.Require().NoError(err)

	queriedImage, err = GetImageByImageUUID(imageUUID)
	suite.Require().NoError(err)
	suite.Require().Nil(queriedImage)
}

func TestImageTestSuite(t *testing.T) {

	defer DB.Disconnect()
	suite.Run(t, new(ImageTestSuite))
}
