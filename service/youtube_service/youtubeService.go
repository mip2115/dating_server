package youtubeservice

import (
	"errors"
	"net/url"
	"sort"
	"strings"

	"code.mine/dating_server/gateway"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/repo"

	"code.mine/dating_server/types"
	"github.com/agnivade/levenshtein"
	stemmer "github.com/agonopol/go-stem"
	"github.com/bbalet/stopwords"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	developerKey = "AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg"
)

// YoutubeController -
type YoutubeController struct {
	gateway gateway.Gateway
	repo    repo.Repo
}

// YOUTUBE ID IS IN THE URL

// GetYoutubeVideoDetails -
func (c *YoutubeController) GetYoutubeVideoDetails(videoURL *string) (*types.UserVideoItem, error) {
	// https://www.googleapis.com/youtube/v3/videos?id=D95qIe5pLuA&key=AIzaSyBhFXscTPZr892Uj5h2wRghkFAqTPYtcEg&part=snippet,statistics,topicDetails
	if videoURL == nil {
		return nil, errors.New("Need url")
	}

	// attempt to get the video is
	u, err := url.Parse(mapping.StrToV(videoURL))
	if err != nil {
		return nil, err
	}
	var videoID *string
	values := u.Query()["id"]
	if len(values) == 0 {
		videoID, err = c.gateway.GetYoutubeVideoID(videoURL)
		if err != nil {
			return nil, errors.New("could not get video ID")
		}
	}
	response, err := c.gateway.GetYoutubeVideoDetails(videoID)
	if err != nil {
		return nil, err
	}

	// baseURL := "https://www.googleapis.com/youtube/v3/videos?"
	// url := fmt.Sprintf("%s&id=%s&key=%s&part=snippet,statistics,topicDetails", baseURL, *videoID, developerKey)
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, err
	// }

	// response := &types.UserVideoItem{}
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// err = json.NewDecoder(resp.Body).Decode(&response)
	// if err != nil {
	// 	return nil, err
	// }
	// if response == nil {
	// 	return nil, errors.New("response is nil")
	// }
	return response, nil
}

// GetYoutubeVideoID -
func (c *YoutubeController) GetYoutubeVideoID(youtubeURL *string) (*string, error) {
	if youtubeURL == nil {
		return nil, errors.New("Need url")
	}
	_, err := c.gateway.GetYoutubeVideoID(youtubeURL)
	if err != nil {
		return nil, err
	}
	return nil, err
	// baseURL := "https://www.googleapis.com/youtube/v3/search?"
	// url := fmt.Sprintf("%spart=%s&maxResults=1&q=%s&type=video&key=%s", baseURL, "snippet", *youtubeURL, developerKey)
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// response := &types.VideoIDResponse{}
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }
	// err = json.NewDecoder(resp.Body).Decode(&response)
	// if err != nil {
	// 	return nil, err
	// }
	// if response == nil {
	// 	return nil, errors.New("response is nil")
	// }
	// if len(response.Items) < 1 {
	// 	return nil, errors.New("response.items is empty")
	// }
	// videoResponse := response.Items[0]
	// if videoResponse.ID == nil {
	// 	return nil, errors.New("item ID is nil")
	// }
	// if videoResponse.ID.VideoID == nil {
	// 	return nil, errors.New("videoID is nil")
	// }
	// return videoResponse.ID.VideoID, nil
}

// https://blog.codecentric.de/en/2017/08/gomock-tutorial/

// GetEligibleUsers -
func (c *YoutubeController) GetEligibleUsers(user *types.User) ([]*types.User, error) {
	// city := user.City
	// partnerGender := user.PartnerGender

	// get all eligible uuids.
	// c, err := DB.GetCollection("users")
	// if err != nil {
	// 	return nil, err
	// }

	filters := &bson.M{
		"zipcode":       user.Zipcode,
		"partnerGender": user.Gender,        // ppl looking for my gender
		"gender":        user.PartnerGender, // their gender is what I want.
	}
	options := options.Find()
	options.SetLimit(int64(50000))
	users, err := c.repo.GetUsersByFilter(filters, options)
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

// RankAndMatchYoutubeVideos -
func (c *YoutubeController) RankAndMatchYoutubeVideos(user *types.User) ([]*types.User, error) {

	if user == nil {
		return nil, errors.New("user is nil")
	}
	// make sure user cant get themselves
	users, err := c.GetEligibleUsers(user)
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
	userVideos, err := c.repo.GetVideosByUserUUID(user.UUID)
	if err != nil {
		return nil, err
	}

	youtubeVideoCandidates, err := c.repo.GetVideosByAllUserUUIDs(userUUIDs)
	if err != nil {
		return nil, err
	}

	sortedVideoScoreList := c.GetSortedVideoList(userVideos, youtubeVideoCandidates)

	processedUUIDs := map[string]bool{}
	sortedUsers := []*types.User{}
	for _, videoScore := range sortedVideoScoreList {
		if !processedUUIDs[*videoScore.Video.UserUUID] {
			processedUUIDs[*videoScore.Video.UserUUID] = true
			sortedUsers = append(sortedUsers, userUUIDToUser[*videoScore.Video.UserUUID])
		}
	}
	if len(sortedUsers) <= 3 {
		return sortedUsers, nil
	}
	return sortedUsers[:3], nil

}

// make sure this is a set
func getTagsFromVideo(video *types.UserVideoItem) []string {
	userTags := []string{}
	for _, tag := range video.Items[0].Snippet.Tags {
		words := strings.Split(tag, " ")
		for _, w := range words {
			userTags = append(userTags, w)
		}
	}
	titleWords := strings.Split(video.Items[0].Snippet.Title, " ")
	for _, w := range titleWords {
		userTags = append(userTags, w)
	}
	descriptionWords := strings.Split(video.Items[0].Snippet.Description, " ")
	for _, w := range descriptionWords {
		userTags = append(userTags, w)
	}
	freeFormText := strings.Join(userTags, " ")
	freeFormText = CleanString(freeFormText)
	userTags = strings.Split(freeFormText, " ")

	processedWords := []string{}
	for _, w := range userTags {
		w = GetStemOfWord(w)
		w = strings.ToLower(w)
		if len(w) > 1 {
			processedWords = append(processedWords, w)
		}
	}
	return processedWords
}

// Score -
type Score struct {
	Score float64
	Video *types.UserVideoItem
}

// we want to do this by video so we can detect a particularly strong match
// between videos
// you can also match by most common words – stem the word
// you can also use this as a dictionary
// https://gist.github.com/dgp/1b24bf2961521bd75d6c
// https://techpostplus.com/youtube-video-categories-list-faqs-and-solutions/#YouTube_video_category_name_and_id_list
// match similar categories

// GetSortedVideoList -
func (c *YoutubeController) GetSortedVideoList(
	userVideos []*types.UserVideoItem,
	candidateVideos []*types.UserVideoItem,
) []Score {

	scores := []Score{}
	for _, userVideo := range userVideos {

		// need to do CleanString somehow on these tags
		userTags := getTagsFromVideo(userVideo)
		userCategoryID := userVideo.Items[0].Snippet.CategoryID
		var totalScore float64

		// for every video in the candidate videos, check how well they match up against
		// the user tags

		for _, video := range candidateVideos {
			videoTags := getTagsFromVideo(video)

			wordFrequencyCandidateVideos := map[string]int{}
			for _, tag := range videoTags {
				wordFrequencyCandidateVideos[tag]++
				// start off with 5 as an arbitrary value
				// remember to clean text of stop words!
				// you can also prob do better here with different levels of word frequnecy
			}
			for _, tag := range userTags {

				if wordFrequencyCandidateVideos[tag] >= 6 {
					totalScore += 6
				} else if wordFrequencyCandidateVideos[tag] >= 4 {
					totalScore += 4
				} else if wordFrequencyCandidateVideos[tag] >= 2 {
					totalScore += 2
				}

			}
			// number of similar words
			similarWordsScore := calculateDistanceScoreBetweenTags(userTags, videoTags)
			totalScore += float64(similarWordsScore) * 2

			categoryScore := calculateCategoryScore(userCategoryID, video.Items[0].Snippet.CategoryID)
			totalScore += categoryScore * 5

			videoScore := Score{
				Score: totalScore,
				Video: video,
			}
			scores = append(scores, videoScore)
		}

	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score < scores[j].Score
	})
	return scores
}

// you could have secondary scores as well
// + scores of if they're both in each others slices vs just 1
func calculateCategoryScore(userCategoryID int, candidateCategoryID int) float64 {
	if VideoCategoryMap[userCategoryID] == nil || VideoCategoryMap[candidateCategoryID] == nil {
		return 0
	}

	containsCandidateCategory := false
	containsUserCategory := false

	for _, category := range VideoCategoryMap[userCategoryID] {
		if category == candidateCategoryID {
			containsCandidateCategory = true
		}
	}
	for _, category := range VideoCategoryMap[candidateCategoryID] {
		if category == userCategoryID {
			containsUserCategory = true
		}
	}

	if containsCandidateCategory && containsUserCategory {
		return 2
	}
	if containsCandidateCategory && containsUserCategory {
		return 1
	}
	return 0
}

// I think you should calculate by word, not by tag
// may want to stem each word in the phrase too
// you can also check if there's like 5 words or something that are in common that are more then 5 chars long
// if so you can give one particular score
// if not, then do the fuzzy matching.
// add in title, description
func calculateDistanceScoreBetweenTags(userTags, videoTags []string) int {
	count := 0
	for _, userTag := range userTags {
		for _, videoTag := range videoTags {

			distance := float64(levenshtein.ComputeDistance(userTag, videoTag))

			// if its not the SAME exact word...
			// but requires less then a few tarnsitions...
			// so we are counting the number of similar words
			if distance != 0 && distance < 4 {
				count++
			}
		}
	}
	// remember that lower the distance, higher the similarity
	return count
}

// CleanString -
func CleanString(text string) string {
	text = strings.Replace(text, "\n", " ", -1)

	// regHTTP := regexp.MustCompile(`/^https.*$/`)
	// regWWW := regexp.MustCompile(`/^www.*$/`)
	// text = reg.ReplaceAllString(text, "")
	if strings.Contains(text, "http") {
		text = ""
	}
	if strings.Contains(text, "https") {
		text = ""
	}
	if strings.Contains(text, "www") {
		text = ""
	}
	text = RemoveStopWords(text)
	return text
}

// RemoveStopWords -
func RemoveStopWords(text string) string {
	cleanedString := stopwords.CleanString(text, "en", false)
	return cleanedString
}

// GetStemsOfText -
func GetStemsOfText(text string) string {
	s := strings.Split(text, " ")
	for i, v := range s {
		s[i] = GetStemOfWord(v)
	}
	return strings.Join(s, " ")
}

// GetStemOfWord –
func GetStemOfWord(word string) string {
	wordAsBytes := []byte((word))
	res := string(stemmer.Stem(wordAsBytes))
	return res
}

// func (youtube *YoutubeController) GetTopMatches(videos)
