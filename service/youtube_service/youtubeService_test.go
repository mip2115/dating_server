package youtubeservice

import (
	"testing"

	"code.mine/dating_server/mapping"

	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type YoutubeServiceTestSuite struct {
	suite.Suite
}

func (suite *YoutubeServiceTestSuite) SetupSuite() {

}
func (suite *YoutubeServiceTestSuite) SetupTest() {

}

func (suite *YoutubeServiceTestSuite) TearDownAllSuite() {

}

func (suite *YoutubeServiceTestSuite) TearDownTest() {
}

func (suite *YoutubeServiceTestSuite) TestGetYoutubeVideoID() {
	query := "https://www.youtube.com/watch?v=D95qIe5pLuA"
	id, err := GetYoutubeVideoID(&query)
	suite.NoError(err)
	suite.NotNil(id)
}

func (suite *YoutubeServiceTestSuite) TestGetYoutubeVideoDetails() {
	query := "https://www.youtube.com/watch?v=D95qIe5pLuA"
	id, err := GetYoutubeVideoID(&query)
	suite.NoError(err)
	suite.NotNil(id)
	response, err := GetYoutubeVideoDetails(id)
	suite.NoError(err)
	suite.NotNil(response)
}

func (suite *YoutubeServiceTestSuite) TestGetEligibleUsers() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)

	youtube := YoutubeController{
		Repo: mockRepo,
	}

	user := &types.User{
		UUID:          mapping.StrToPtr("some-uuid"),
		PartnerGender: mapping.StrToPtr("partnerGender"),
		City:          mapping.StrToPtr("city"),
	}
	returnedUsers := []*types.User{user}

	mockRepo.EXPECT().GetUsersByFilter(gomock.Any(), gomock.Any()).Return(returnedUsers, nil)
	users, err := youtube.GetEligibleUsers(user)
	suite.Require().NoError(err)
	suite.Require().NotNil(users)
	// suite.Require().Equal(1, len(users))
}

func (suite *YoutubeServiceTestSuite) TestRankAndMatchYoutubeVideos() {
	candidateVideos := []*types.UserVideoItem{
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-1"),
			UUID:     mapping.StrToPtr("uuid-1"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("tSbM2r5FJHc"),
					Snippet: types.Snippet{
						Title: "Remove All Negative Blockages Wipe Out Subconscious Negative Patterns, 432 Hz Boost Your Aura",
						Tags: []string{
							"Remove All Negative Blockages",
							"Erase Subconscious Negative Patterns",
							"432 Hz",
							"boost your aura",
							"432 Hz Boost Your Aura",
							"remove negative energy",
							"remove mental blockages",
							"remove negative thoughts",
							"remove negative blocks",
							"let go of mental blockages",
							"remove negativity",
							"mental blockages",
							"dissolve negative patterns",
							"remove all negative blockages",
							"jason stephenson",
							"chakra healing music",
						},
						CategoryID: 10,
					},
				},
			},
		},
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-2"),
			UUID:     mapping.StrToPtr("uuid-2"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("0Rc7l2DgSpc"),
					Snippet: types.Snippet{
						Title: "Adventure Time | Simon & Marcy | Cartoon Network",
						Tags: []string{
							"adventure time",
							"fin and jack",
							"adventure time cartoon",
							"new adventure time",
							"adventure time new episodes 2015",
							"adventure time new episode",
							"adventure time last episode",
							"avenger time",
							"adventure time new episodes",
							"finn and jake",
							"ice king",
							"fin and jake",
							"Jake",
							"Ice King",
							"Marceline",
							"adventure time episode",
							"adventure time full",
							"the adventure time",
							"The Snow Golem's New Friend",
							"Thanksgiving",
							"adventure",
							"time",
							"cartoon network",
							"cartoon network adventure time",
							"cartoon",
						},
						CategoryID: 1,
					},
				},
			},
		},
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-3"),
			UUID:     mapping.StrToPtr("uuid-3"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("9_LTYgUL8t8"),
					Snippet: types.Snippet{
						Title: "Rick and Morty - Best Moments | Season 3",
						Tags: []string{
							"Rick and Morty",
							"Morty",
							"Rick",
							"Best Moments",
							"Season 3",
							"Series",
							"Funny",
							"Funny moments",
						},
						CategoryID: 1,
					},
				},
			},
		},
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-4"),
			UUID:     mapping.StrToPtr("uuid-4"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("VxzLYJa87Bk"),
					Snippet: types.Snippet{
						Title: "Elon Musk Shares New Starship Details | SpaceX in the News",
						Tags: []string{
							"spacex starship",
							"spacex explosion",
							"spacex starship test",
							"spacex starship sn8",
							"spacex landing",
							"spacex starship landing",
							"spacex starship flight",
							"spacex starship launch",
							"spacex news",
							"spacex starship news",
							"spacex starlink",
							"elon musk",
							"spacex dragon",
						},
						CategoryID: 28,
					},
				},
			},
		},
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-5"),
			UUID:     mapping.StrToPtr("uuid-5"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("q22uHBl9qxw"),
					Snippet: types.Snippet{
						Title: "Solitude ● lofi hip hop mix",
						Tags: []string{
							"lofi hip hop",
							"chillhop",
							"lofi",
							"soulful",
							"chill beats",
							"라디오",
							"로파이힙합",
							"lofi radio",
							"jazz hip hop",
							"chill study beats",
							"chill gaming beats",
							"chill relax beats",
							"beats to study to",
							"anime beats",
							"chillhop records",
							"chill mix",
							"study mix",
							"lo-fi beats",
							"lo fi hip hop",
							"nujabes",
							"chillhop raccoon",
							"best lofi hip hop",
							"dreamy",
						},
						CategoryID: 10,
					},
				},
			},
		},
		&types.UserVideoItem{
			UserUUID: mapping.StrToPtr("user-uuid-6"),
			UUID:     mapping.StrToPtr("uuid-6"),
			Items: []types.DetailsItem{
				types.DetailsItem{
					ID: mapping.StrToPtr("ltwzbOQJrcM"),
					Snippet: types.Snippet{
						Title:       "Best Makeup Transformations 2020 | New Makeup Tutorials Compilation",
						Description: "Best Makeup Transformations 2020 | New Makeup Tutorials Compilation \n\nCheck Out These Amazing Makeup Artists:\nhttps://www.instagram.com/huberbeauty/\nhttps://www.instagram.com/ccclarke/\nhttps://www.instagram.com/alxcext/\nhttps://www.instagram.com/lenkalul/\nhttps://www.instagram.com/itsisabelbedoya/\nhttps://www.instagram.com/noazitoun_makeupartist/\nhttps://www.instagram.com/adi_katzanelson_makeupartist/\nhttps://www.instagram.com/inbal_eitan/\nhttps://www.instagram.com/xthuyle/\nhttps://www.instagram.com/mayamua__/\nhttps://www.instagram.com/glambby_/\nhttps://www.instagram.com/monserrathsglam/\nhttps://www.instagram.com/b.q.a.j/\nhttps://www.instagram.com/bexcxmpbell/\nhttps://www.instagram.com/ortal_azizada/\nhttps://www.instagram.com/lenaglams/",
						Tags:        []string{},
						CategoryID:  26,
					},
				},
			},
		},
	}

}

func TestYoutubeServiceTestSuite(t *testing.T) {

	suite.Run(t, new(YoutubeServiceTestSuite))
}
