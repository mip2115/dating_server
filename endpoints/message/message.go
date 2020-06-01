package message

import (
	"encoding/json"
	"github.com/kama/server/service/aws_service"
	"github.com/kama/server/service/image_service"
	"github.com/kama/server/service/message_service"
	"github.com/kama/server/service/user_service"
	"github.com/kama/server/types"
	"net/http"
)

var (
	nPerPage = 4
)

// create middleware that will ensure
// only the sender of this is in the from position
// also you need to use pagination and return in sorrted order by date
// use aggregate?
func AddMessage(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	messageRequest := &types.MessageRequest{}
	err := json.NewDecoder(r.Body).Decode(messageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := message_service.AddMessage(messageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msgResponse := &types.MessageResponse{}
	msgResponse.Message = msg
	msgResponse.Token = &tkString
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(msgResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	tkString := r.Context().Value("tokenString").(string)
	mr := &types.MessageRequest{}
	err := json.NewDecoder(r.Body).Decode(mr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msgs, err := message_service.GetMessages(mr, nPerPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.MessagesObjectResponse{}
	response.Token = &tkString
	response.Messages = msgs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
