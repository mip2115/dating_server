package aws_service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/kama/server/DB"
	AWSSetup "github.com/kama/server/DB/aws"
	//us "github.com/kama/server/service/user_service"
	//"github.com/kama/server/types"
	//AWSSetup "github.com/code/kama/server/DB/aws"
	//image "github.com/kama/server/service/image_service"
	"github.com/nu7hatch/gouuid"
	"os"
)

///Users/M/go/src/github.com/code/kama/server/service/aws_service/aws_service.go
var (
	imageFolder  = "user-images"
	linkToBucket = "https://kama-documents-public.s3.amazonaws.com/"
)

func UploadImageToS3(fileBytes []byte, userUUID *string) (*string, *string, error) {
	decodedString, err := base64.StdEncoding.DecodeString(string(fileBytes))
	if err != nil {
		return nil, nil, err
	}
	s, err := AWSSetup.GetSession()
	if err != nil {
		return nil, nil, err
	}
	buf := bytes.NewBuffer([]byte(decodedString))
	size := buf.Len()
	buffer := buf.Bytes()

	newFileNameUUID, err := uuid.NewV4()
	if err != nil {
		return nil, nil, err
	}

	fileNameForS3 := fmt.Sprintf("%s/%s/%s.jpg", imageFolder, *userUUID, newFileNameUUID.String())

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("BUCKET_NAME")),
		Key:                  aws.String(fileNameForS3),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String("image/jpeg"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return nil, nil, err
	}

	linkToAdd := fmt.Sprintf("%s%s", linkToBucket, fileNameForS3)
	/*
		img, err := image.CreateImage(userUUID, &linkToAdd, &fileNameForS3)
		if err != nil {
			return err
		}

		/*
		//linkToAdd := fmt.Sprintf("%s%s", linkToBucket, fileNameForS3)
		err = us.SaveUserImage(userUUID, img.UserUUID)
		if err != nil {
			return err
		}
	*/
	// after we've uploaded, add the image to the users thing
	return &linkToAdd, &fileNameForS3, nil
}

func DeleteImageFromS3(key *string) error {
	s, err := AWSSetup.GetSession()
	if err != nil {
		return err
	}
	_, err = s3.New(s).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(*key),
	})
	if err != nil {
		return err
	}
	return nil
}
