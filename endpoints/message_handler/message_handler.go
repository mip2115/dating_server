package message_handler

import (
	"encoding/json"
	"net/http"

	"code.mine/dating_server/service/message_service"
	"code.mine/dating_server/types"
)

var (
	nPerPage = 4
)

// create middleware that will ensure
// only the sender of this is in the from position
// also you need to use pagination and return in sorrted order by date
// use aggregate?

// AddMessage –
func AddMessage(w http.ResponseWriter, r *http.Request) {
	//userUUID := r.Context().Value("userUUID").(string)
	//tkString := r.Context().Value("tokenString").(string)
	messageRequest := &types.MessageRequest{}
	err := json.NewDecoder(r.Body).Decode(messageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = message_service.AddMessage(messageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// msgResponse := &types.MessageResponse{}
	// msgResponse.Message = interface{}
	// msgResponse.Token = &tkString
	msgResponse := map[string]interface{}{}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(msgResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// GetMessages –
func GetMessages(w http.ResponseWriter, r *http.Request) {
	// tkString := r.Context().Value("tokenString").(string)
	// mr := &types.MessageRequest{}
	// err := json.NewDecoder(r.Body).Decode(mr)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// msgs, err := message_service.GetMessages(mr, nPerPage)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	response := map[string]interface{}{}
	// response.Token = &tkString
	// response.Messages = msgs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
