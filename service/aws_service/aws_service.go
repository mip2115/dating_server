package awsservice

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"code.mine/dating_server/gateway"
	"code.mine/dating_server/mapping"
	uuid "github.com/satori/go.uuid"
)

// AWSController -
type AWSController struct {
	gateway gateway.Gateway
}

///Users/M/go/src/github.com/code/kama/server/service/aws_service/aws_service.go
var (
	imageFolder  = "user-images"
	linkToBucket = "https://kama-documents-public.s3.amazonaws.com/"
)

// UploadImageToS3 -
func (c *AWSController) UploadImageToS3(fileBytes []byte, userUUID *string) (*string, error) {
	decodedString, err := base64.StdEncoding.DecodeString(string(fileBytes))
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte(decodedString))
	size := int64(buf.Len())
	buffer := buf.Bytes()
	newFileNameUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	fileNameForS3 := fmt.Sprintf("%s/%s/%s.jpg", imageFolder, mapping.StrToV(userUUID), newFileNameUUID.String())
	err = c.gateway.UploadUserImageToS3(fileNameForS3, buffer, size)
	if err != nil {
		return nil, err
	}

	urlForUploadedImage := fmt.Sprintf("%s%s", linkToBucket, fileNameForS3)
	return &urlForUploadedImage, nil
}

// DeleteImageFromS3 -
func (c *AWSController) DeleteImageFromS3(key *string) error {
	err := c.gateway.DeleteUserImageFromS3(key)
	if err != nil {
		return err
	}
	return nil
}
