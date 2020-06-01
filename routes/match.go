package routes

import (
	"../auth"
	"../endpoints/match"
	"net/http"
)

func GetDeleteMatchHandler() http.Handler {
	handler := http.HandlerFunc(match.DeleteMatch)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}
