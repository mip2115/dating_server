package search_endpoint

import (
	"errors"
	"github.com/kama/server/service/search_service"
	"net/http"
)

func RetrieveUsers(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	searchRequest := &types.SearchRequest{}
	err := json.NewDecoder(r.Body).Decode(searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if searchRequest.SkipValue == nil {
		http.Error(w, errors.New("need to provide a skip value"), http.StatusBadRequest)
		return
	}
	if searchRequest.UserUUID == nil {
		http.Error(w, errors.New("need to provide a user uuid"), http.StatusBadRequest)
		return
	}
	skipValue := *searchRequest.SkipValue
	userUUID := searchRequest.UserUUID
	topCalculatedUsers, err := search_service.CalculateTopUsers(userUUID, skipValue)

	// get list of top 3 users
	// topUsers = service.CalculateTopUsers(userID)

	// return topUsers
}
