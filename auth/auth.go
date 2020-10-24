package auth

import (
	"code.mine/dating_server/types"
	"github.com/dgrijalva/jwt-go"

	//"golang.org/x/crypto/bcrypt"
	"context"
	//"errors"
	"net/http"
	"strings"
	"time"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// GenerateJWT -
func GenerateJWT(uuid *string, expiresAt time.Time) (string, error) {
	//expiresAt := time.Now().Add(3 * time.Minute)
	claims := &types.Token{
		UserUUID: *uuid,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}

	/*
		tk := &types.Token{
			UserID: ID,
			StandardClaims: &jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}
	*/

	return tokenString, nil
}

// https://www.sohamkamani.com/golang/2019-01-01-jwt-authentication/
func VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token") //Grab the token from the header
		header = strings.TrimSpace(header)

		if header == "" {
			http.Error(w, "Missing header", http.StatusBadRequest)
			return
		}
		claims := types.Token{}

		tkn, err := jwt.ParseWithClaims(header, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "token not valid", http.StatusBadRequest)
			return
		}

		/*
			v := map[string]string{
				"userID": claims.UserID,
			}
		*/

		ctx := context.WithValue(r.Context(), "userUUID", tkn.Claims.(*types.Token).UserUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// https://stackoverflow.com/questions/51201056/testing-golang-middleware-that-modifies-the-request
func RefreshJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var tknStr = r.Header.Get("x-access-token") //Grab the token from the header
		tknStr = strings.TrimSpace(tknStr)

		if tknStr == "" {
			http.Error(w, "Missing header", http.StatusBadRequest)
			return
		}
		claims := types.Token{}

		tkn, err := jwt.ParseWithClaims(tknStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "token not valid", http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(3 * time.Minute)
		claims.StandardClaims = jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		/*
			v := map[string]string{
				"userID": claims.UserID,
			}
		*/

		ctx := context.WithValue(r.Context(), "userUUID", tkn.Claims.(*types.Token).UserUUID)
		ctx = context.WithValue(ctx, "tokenString", tokenString)
		//context.Set(r, "decoded", token.Claims)
		//next(w, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
