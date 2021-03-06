package match_handler

import (
	"encoding/json"
	//"fmt"
	// "../../mapping"
	"net/http"
	// "../../DB"
	// "../../auth"
	"code.mine/dating_server/service/match_service"
	"code.mine/dating_server/types"

	//"context"
	"github.com/gorilla/mux"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"golang.org/x/crypto/bcrypt"
	//"strings"
	//"bytes"
)

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	//userID := r.Context().Value("userID").(string)
	tkString := r.Context().Value("tokenString").(string)
	m := &types.Match{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := match_service.CreateMatch(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m.UUID = res
	response := types.MatchResponse{}
	response.Match = m
	response.Token = &tkString
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// if they get blocked or something
// also make sure that only the peopel int he match can delete it
func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	params := mux.Vars(r)
	deleteMatchID := params["mid"]
	profileBID := params["bid"]
	m := &types.MatchRequest{}
	m.UserOneUUID = &userID
	m.UserTwoUUID = &profileBID
	m.UUID = &deleteMatchID
	err := match_service.DeleteMatch(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.MatchResponse{}
	response.Token = &tkString
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetMatchMessages(w http.ResponseWriter, r *http.Request) {

}
