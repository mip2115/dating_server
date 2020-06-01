package routes

import (
	"../auth"
	"../endpoints/message"
	"net/http"
)

func GetAddMessageHandler() http.Handler {

	handler := http.HandlerFunc(message.AddMessage)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h

}
