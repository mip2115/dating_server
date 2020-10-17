package imageservice

import (
	"errors"

	"code.mine/dating_server/mapping"
	"code.mine/dating_server/repo"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
)

// ImageController -
type ImageController struct {
	repo repo.Repo
}

// TODO – figure out more graceful way to deal with rank
// deal with images that already exist as well
func (c *ImageController) CreateImage(imageToUpload *types.Image) (*types.Image, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	imageToUpload.UUID = mapping.StrToPtr(uuid.String())
	// if the current rank exists then just delete that
	err = c.repo.CreateImage(imageToUpload)
	if err != nil {
		return nil, err
	}
	return imageToUpload, nil
}

// GetImagesByUserUUID -
func (c *ImageController) GetImagesByUserUUID(userUUID *string) ([]*types.Image, error) {
	if userUUID == nil {
		return nil, errors.New("need user uuid")
	}
	results, err := c.repo.GetImagesByUserUUID(userUUID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GetImageByImageUUID -
func (c *ImageController) GetImageByImageUUID(uuid *string) (*types.Image, error) {
	if uuid == nil {
		return nil, errors.New("need image uuid")
	}
	img, err := c.repo.GetImageByImageUUID(uuid)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// DeleteImage -
func (c *ImageController) DeleteImage(imageUUID *string) error {
	if imageUUID == nil {
		return errors.New("need image uuid")
	}
	err := c.repo.DeleteImage(imageUUID)
	if err != nil {
		return err
	}
	return nil
}

