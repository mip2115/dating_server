package gateway

//go:generate mockgen -destination=mocks/gateway/mock_gateway.go -package=gatewaymocks code.mine/dating_server/gateway Gateway

import "code.mine/dating_server/types"

type Gateway interface {
	DeleteUserImageFromS3(key *string) error
	GetYoutubeVideoDetails(videoID *string) (*types.UserVideoItem, error)
	GetYoutubeVideoID(youtubeURL *string) (*string, error)
	UploadUserImageToS3(fileNameForS3 string, body []byte, size int64) error
}
