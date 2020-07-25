package routes

import (
	"net/http"

	"code.mine/dating_server/auth"
	"code.mine/dating_server/endpoints/message_handler"
)

func GetAddMessageHandler() http.Handler {

	handler := http.HandlerFunc(message_handler.AddMessage)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h

}
