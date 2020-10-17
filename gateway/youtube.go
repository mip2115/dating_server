package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"code.mine/dating_server/mapping"

	"code.mine/dating_server/types"
)

const (
	developerKey = "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
)

// GetYoutubeVideoDetails -
func GetYoutubeVideoDetails(videoID *string) (*types.UserVideoItem, error) {
	baseURL := "https://www.googleapis.com/youtube/v3/videos?"
	url := fmt.Sprintf("%s&id=%s&key=%s&part=snippet,statistics,topicDetails", baseURL, mapping.StrToV(videoID), developerKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response := &types.UserVideoItem{}
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

func GetYoutubeVideoID(youtubeURL *string) (*string, error) {
	baseURL := "https://www.googleapis.com/youtube/v3/search?"
	url := fmt.Sprintf("%spart=%s&maxResults=1&q=%s&type=video&key=%s", baseURL, "snippet", mapping.StrToV(youtubeURL), developerKey)
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
