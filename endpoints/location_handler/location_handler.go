package location_handler

import (
	"encoding/json"
	//"fmt"
	// "../../mapping"
	"net/http"
	// "../../DB"
	// "../../auth"
	"github.com/kama/server/service/location_service"
	"github.com/kama/server/types"
	//"context"
	"github.com/gorilla/mux"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"golang.org/x/crypto/bcrypt"
	//"strings"
	//"bytes"
)

func CreateLocation(w http.ResponseWriter, r *http.Request) {
	//userID := r.Context().Value("userID").(string)
	tkString := r.Context().Value("tokenString").(string)
	loc := &types.Location{}
	err := json.NewDecoder(r.Body).Decode(loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := location_service.CreateLocation(loc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	m.MatchUUID = res
	response := types.MatchResponse{}
	response.Match = m
	response.Token = &tkString
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
