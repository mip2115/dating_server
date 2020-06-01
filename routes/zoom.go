package routes

import (
	"../endpoints/zoom"
	"net/http"
)

func GetZoomHandler() http.Handler {
	handler := http.HandlerFunc(zoom.ZoomCallback)
	return handler
}

func GetZoomTokenHandler() http.Handler {
	handler := http.HandlerFunc(zoom.ZoomTokenCallback)
	return handler
}
