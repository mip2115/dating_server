package youtubeService

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"code.mine/dating_server/types"
)

const (
	developerKey = "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
)

func GetYoutubeVideo(youtubeURL string) {
	// developerKey := "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
	// baseURL := "https://www.googleapis.com/youtube/v3/search?"

}

// GetYoutubeVideoDetails -
func GetYoutubeVideoDetails(videoID *string) (*types.VideoDetailsResponse, error) {
	// https://www.googleapis.com/youtube/v3/videos?id=D95qIe5pLuA&key=AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg&part=snippet,statistics,topicDetails
	if videoID == nil {
		return nil, errors.New("Need url")
	}
	baseURL := "https://www.googleapis.com/youtube/v3/videos?"
	url := fmt.Sprintf("%s&id=%s&key=%s&part=snippet,statistics,topicDetails", baseURL, *videoID, developerKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response := &types.VideoDetailsResponse{}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response is nil")
	}
	return response, nil
}

// GetYoutubeVideoID -
func GetYoutubeVideoID(youtubeURL *string) (*string, error) {
	if youtubeURL == nil {
		return nil, errors.New("Need url")
	}
	baseURL := "https://www.googleapis.com/youtube/v3/search?"
	url := fmt.Sprintf("%spart=%s&maxResults=1&q=%s&type=video&key=%s", baseURL, "snippet", *youtubeURL, developerKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response := &types.VideoIDResponse{}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response is nil")
	}
	if len(response.Items) < 1 {
		return nil, errors.New("response.items is empty")
	}
	videoResponse := response.Items[0]
	if videoResponse.ID == nil {
		return nil, errors.New("item ID is nil")
	}
	if videoResponse.ID.VideoID == nil {
		return nil, errors.New("videoID is nil")
	}
	return videoResponse.ID.VideoID, nil

}
