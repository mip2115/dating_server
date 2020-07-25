package routes

import (
	"net/http"

	"code.mine/dating_server/auth"
	"code.mine/dating_server/endpoints/match_handler"
)

func GetDeleteMatchHandler() http.Handler {
	handler := http.HandlerFunc(match_handler.DeleteMatch)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}
