package gateway

import (
	"bytes"
	"os"

	AWSSetup "code.mine/dating_server/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadUserImageToS3 -
func UploadUserImageToS3(fileNameForS3 string, buffer []byte, size int64) error {
	s, err := AWSSetup.GetSession()
	if err != nil {
		return err
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("BUCKET_NAME")),
		Key:                  aws.String(fileNameForS3),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String("image/jpeg"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserImageFromS3 -
func DeleteUserImageFromS3(key *string) error {
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
