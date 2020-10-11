package youtubeservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"

	"code.mine/dating_server/repo"
	"code.mine/dating_server/types"
	"github.com/agnivade/levenshtein"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	developerKey = "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
)

type YoutubeController struct {
	Repo repo.Repo
}

func GetYoutubeVideo(youtubeURL string) {
	// developerKey := "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
	// baseURL := "https://www.googleapis.com/youtube/v3/search?"

}

// GetYoutubeVideoDetails -
func GetYoutubeVideoDetails(videoID *string) (*types.UserVideoItem, error) {
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

// https://blog.codecentric.de/en/2017/08/gomock-tutorial/
// GetEligibleUsers -
func (youtube *YoutubeController) GetEligibleUsers(user *types.User) ([]*types.User, error) {
	city := user.City
	partnerGender := user.PartnerGender

	// get all eligible uuids.
	// c, err := DB.GetCollection("users")
	// if err != nil {
	// 	return nil, err
	// }

	filters := bson.M{
		"city":          *city,
		"partnerGender": *partnerGender,
	}
	options := options.Find()
	options.SetLimit(int64(50000))
	users, err := youtube.Repo.GetUsersByFilter(&filters, options)
	if err != nil {
		return nil, err
	}

	// cursor, err := c.Find(context.Background(), filters, options)
	// users := []*types.User{}
	// if err = cursor.All(context.Background(), &users); err != nil {
	// 	return nil, err
	// }
	return users, nil
}

// TODO – add in title words to tags
func (youtube *YoutubeController) RankAndMatchYoutubeVideos(user *types.User) ([]*types.User, error) {

	if user == nil {
		return nil, errors.New("user is nil")
	}
	// make sure user cant get themselves
	users, err := youtube.GetEligibleUsers(user)
	if err != nil {
		return nil, err
	}

	userUUIDToUser := map[string]*types.User{}
	userUUIDs := []*string{}
	for _, user := range users {
		userUUIDs = append(userUUIDs, user.UUID)
		userUUIDToUser[*user.UUID] = user
	}

	// get video by user
	userVideos, err := youtube.Repo.GetVideosByUserUUID(user.UUID)
	if err != nil {
		return nil, err
	}

	userTags := getTagsFromVideos(userVideos)

	youtubeVideoCandidates, err := youtube.Repo.GetVideosByAllUserUUIDs(userUUIDs)
	if err != nil {
		return nil, err
	}

	sortedVideoScoreList := youtube.GetSortedVideoList(userTags, youtubeVideoCandidates)

	processedUUIDs := map[string]bool{}
	sortedUsers := []*types.User{}
	for _, videoScore := range sortedVideoScoreList {
		if !processedUUIDs[*videoScore.Video.UserUUID] {
			processedUUIDs[*videoScore.Video.UserUUID] = true
			sortedUsers = append(sortedUsers, userUUIDToUser[*videoScore.Video.UserUUID])
		}
	}
	return sortedUsers, nil

}

// make sure this is a set
func getTagsFromVideos(videos []*types.UserVideoItem) []string {
	userTags := []string{}
	for _, video := range videos {
		for _, tag := range video.Items[0].Snippet.Tags {
			userTags = append(userTags, tag)
		}
	}
	return userTags

}

type Score struct {
	Score float64
	Video *types.UserVideoItem
}

// we want to do this by video so we can detect a particularly strong match
// between videos
func (youtube *YoutubeController) GetSortedVideoList(userTags []string, candidateVideos []*types.UserVideoItem) []Score {

	// for every video in the candidate videos, check how well they match up against
	// the user tags
	scores := []Score{}
	for _, video := range candidateVideos {
		videoTags := getTagsFromVideos([]*types.UserVideoItem{video})
		score := calculateDistanceScoreBetweenTags(userTags, videoTags)
		videoScore := Score{
			Score: score,
			Video: video,
		}
		scores = append(scores, videoScore)
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})
	return scores
}

// may want to stem each word in the phrase too
func calculateDistanceScoreBetweenTags(userTags, videoTags []string) float64 {
	var totalScore float64

	lengthOfLongest := math.Max(float64(len(userTags)), float64(len(videoTags)))
	for _, userTag := range userTags {
		for _, videoTag := range videoTags {
			distance := float64(levenshtein.ComputeDistance(userTag, videoTag))
			totalScore += distance
		}
	}

	// remember that lower the distance, higher the similarity
	return 1 - (totalScore / lengthOfLongest)
}

// func (youtube *YoutubeController) GetTopMatches(videos)
