package functional

import (
	"../../DB"
	"../../types"
	"../urls"
	"bytes"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"net/http"
	//urlpkg "net/url"
	"testing"
)

var (
	user1Email     = "testing1@gmail.com"
	user1Password  = "testing123"
	user1FirstName = "user1firstname"
	user1dob       = "dob"
	user1mobile    = "mobile"

	user2Email     = "testing2@gmail.com"
	user2Password  = "testing123"
	user2FirstName = "user2firstname"
	user2dob       = "dob"
	user2mobile    = "mobile"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type UserTestSuite struct {
	suite.Suite
}

func (suite *UserTestSuite) SetupTest() {

}

func (suite *UserTestSuite) TearDownAllSuite() {
	err := DB.Disconnect()
	suite.Nil(err)

}

func (suite *UserTestSuite) TearDownTest() {

	c, err := DB.GetCollection("users")
	suite.Nil(err)
	err = c.Drop(context.Background())
	suite.Nil(err)

	c, err = DB.GetCollection("matches")
	suite.Nil(err)
	err = c.Drop(context.Background())
	suite.Nil(err)
}

func (suite *UserTestSuite) TestCreateUser() {

	userResponse, err := createUser(user1Email, user1Password, user1FirstName, user1dob, user1mobile)
	suite.Nil(err)

	_, err = getUser(userResponse.User.Email)
	suite.Nil(err)
}

func (suite *UserTestSuite) TestUpdateUser() {
	_, err := createUser(user1Email, user1Password, user1FirstName, user1dob, user1mobile)
	suite.Nil(err)

	ur, err := loginUser(user1Email, user1Password)
	suite.Nil(err)

	url := urls.UPDATE_USER
	ct := "application/json"
	rb, err := json.Marshal(map[string]string{
		"gender": "DJ",
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rb))
	suite.Nil(err)
	req.Header.Set("x-access-token", *ur.Token)
	req.Header.Set("Content-Type", ct)

	client := http.Client{}
	resp, err := client.Do(req)
	suite.Nil(err)
	defer resp.Body.Close()

	//result := make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)

	//err = json.Unmarshal(body, &result)
	//suite.Nil(err)

	userResponse := types.UserResponse{}

	err = json.Unmarshal(body, &userResponse)
	suite.Nil(err)

	suite.NotNil(userResponse.Token)
	suite.NotNil(userResponse.User)

	suite.Equal("DJ", *userResponse.User.Gender)

}

func (suite *UserTestSuite) TestLoginUser() {
	_, err := createUser(user1Email, user1Password, user1FirstName, user1dob, user1mobile)
	suite.Nil(err)

	ur, err := loginUser(user1Email, user1Password)
	suite.Nil(err)

	suite.Nil(err)
	suite.NotNil(ur.Token)
	suite.NotNil(ur.User)

}

func (suite *UserTestSuite) TestDeleteUser() {
	_, err := createUser(user1Email, user1Password, user1FirstName, user1dob, user1mobile)
	suite.Nil(err)

	ur, err := loginUser(user1Email, user1Password)
	suite.Nil(err)
	token := *ur.Token

	url := urls.DELETE_USER
	ct := "application/json"
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("x-access-token", token)
	req.Header.Set("Content-Type", ct)

	client := http.Client{}
	resp, err := client.Do(req)
	suite.Nil(err)
	defer resp.Body.Close()

	result := make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	suite.Nil(err)

	err = json.Unmarshal(body, &result)
	suite.Nil(err)

	userResponse := types.UserResponse{}

	err = json.Unmarshal(body, &userResponse)
	suite.Nil(err)

	suite.True(*userResponse.Op)

	_, err = getUser(ur.User.Email)
	suite.NotNil(err)

}

func (suite *UserTestSuite) TestLikeUser() {
	_, err := createUser(user1Email, user1Password, user1FirstName, user1dob, user1mobile)
	suite.Nil(err)

	_, err = createUser(user2Email, user2Password, user2FirstName, user2dob, user2mobile)
	suite.Nil(err)

	ur1, err := loginUser(user1Email, user1Password)
	suite.Nil(err)
	user1 := ur1.User

	ur2, err := loginUser(user2Email, user2Password)
	suite.Nil(err)
	user2 := ur2.User

	// user1 likes user2
	mr, err := likeUser(*user2.UserID, *ur1.Token)
	suite.Nil(err)

	// make sure user2's list of people who liked them has user1
	// but user1 doesn't have user 2
	user1, err = getUser(user1.Email)
	suite.Nil(err)

	user2, err = getUser(user2.Email)
	suite.Nil(err)

	suite.Nil(mr.Match)

	u1UsersLikedMe := *user1.UsersLikedMe
	u1UserMatches := *user1.Matches
	u2UsersLikedMe := *user2.UsersLikedMe
	u2UserMatches := *user2.Matches

	suite.Empty(u1UsersLikedMe)
	suite.Empty(u1UserMatches)
	suite.Empty(u2UserMatches)
	suite.Equal(1, len(u2UsersLikedMe))
	suite.Equal(0, len(u1UsersLikedMe))

	mr, err = likeUser(*user1.UserID, *ur2.Token)
	suite.Nil(err)

	user1, err = getUser(user1.Email)
	suite.Nil(err)

	user2, err = getUser(user2.Email)
	suite.Nil(err)

	suite.NotNil(mr.Match)

	u1UsersLikedMe = *user1.UsersLikedMe
	u1UserMatches = *user1.Matches
	u2UsersLikedMe = *user2.UsersLikedMe
	u2UserMatches = *user2.Matches

	suite.Empty(u1UsersLikedMe)
	suite.Empty(u2UsersLikedMe)
	suite.NotEmpty(u1UserMatches)
	suite.NotEmpty(u2UserMatches)

}

func likeUser(likedID string, token string) (*types.MatchResponse, error) {

	url := urls.LIKE_USER + "/" + likedID
	ct := "application/json"
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("x-access-token", token)
	req.Header.Set("Content-Type", ct)

	//q := req.URL.Query()
	//q.Add("id", likedID)
	//req.URL.RawQuery = q.Encode()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	st := string(body)
	_ = st

	matchResponse := types.MatchResponse{}
	err = json.Unmarshal(body, &matchResponse)
	if err != nil {
		return nil, err
	}
	return &matchResponse, nil

}

func loginUser(email string, password string) (*types.UserResponse, error) {
	rb, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	url := urls.LOGIN_USER
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	userResponse := types.UserResponse{}

	result := make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}
	return &userResponse, nil

}

func createUser(email string,
	password string, firstName string,
	dob string, mobile string) (*types.UserResponse, error) {

	url := urls.CREATE_USER
	ct := "application/json"

	rb, err := json.Marshal(map[string]string{
		"email":      email,
		"password":   password,
		"first_name": firstName,
		"dob":        dob,
		"mobile":     mobile,
	})
	resp, err := http.Post(url, ct, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//result := make(map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	userResponse := types.UserResponse{}

	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}
	return &userResponse, nil

}

func getUser(email *string) (*types.User, error) {
	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}

	u := types.User{}
	res := c.FindOne(context.Background(), bson.M{"email": email})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(&u)
	return &u, nil
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run

func TestUserTestSuite(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	_, err := DB.SetupDB()
	if err != nil {
		panic(err.Error())
	}

	suite.Run(t, new(UserTestSuite))

}
