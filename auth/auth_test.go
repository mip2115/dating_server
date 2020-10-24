package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func (suite *AuthTestSuite) SetupSuite() {

}
func (suite *AuthTestSuite) SetupTest() {

}

func (suite *AuthTestSuite) TearDownAllSuite() {

}

func (suite *AuthTestSuite) TearDownTest() {
}

func (suite *AuthTestSuite) TestGenerateJWT() {
	expiresAt := time.Now().Add(5 * time.Minute)
	uuid := mapping.StrToPtr("some-uuid")
	token, err := GenerateJWT(uuid, expiresAt)
	suite.Require().NoError(err)
	suite.Require().NotNil(token)

	claims := types.Token{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	// expiration happens here
	suite.Require().NoError(err)
	suite.Require().NotNil(parsedToken)
	userUUID := parsedToken.Claims.(*types.Token).UserUUID
	suite.Require().Equal(userUUID, mapping.StrToV(uuid))

	ett := time.Now().Add(-3 * time.Minute)
	token, err = GenerateJWT(uuid, ett)
	suite.Require().NoError(err)
	suite.Require().NotNil(token)

	claims = types.Token{}
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	suite.Require().Error(err)
}

func (suite *AuthTestSuite) TestVerifyJWT() {
	uuid := "uuid"
	expiresAt := time.Now().Add(5 * time.Minute)
	tokenForHeader, err := GenerateJWT(&uuid, expiresAt)
	suite.Require().NoError(err)
	fn1 := func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value("userUUID").(string)
		suite.Require().Equal(uuid, userUUID)
	}
	w := httptest.NewRecorder()

	h1 := VerifyJWT(http.HandlerFunc(fn1))
	req, _ := http.NewRequest("GET", "https://example.com/foo/", nil)
	req.Header.Set("x-access-token", tokenForHeader)
	h1.ServeHTTP(w, req)

	expiresAt = time.Now().Add(-5 * time.Minute)
	tokenForHeader, err = GenerateJWT(&uuid, expiresAt)
	suite.Require().NoError(err)
	fn1 = func(w http.ResponseWriter, r *http.Request) {
		// this should fail if we can pass the middleware.
		// if it doesn't fail, we know the middleware stopped the req
		suite.Require().Equal(true, false)
	}
	w = httptest.NewRecorder()

	h1 = VerifyJWT(http.HandlerFunc(fn1))
	req, _ = http.NewRequest("GET", "https://example.com/foo/", nil)
	req.Header.Set("x-access-token", tokenForHeader)
	h1.ServeHTTP(w, req)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
