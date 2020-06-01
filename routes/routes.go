package routes

import (
	//"../auth"

	//"../services/user"
	"github.com/gorilla/mux"
	//"net/http"
)

// auth with google // https://dev.to/douglasmakey/oauth2-example-with-go-3n8a

func CreateRoutes(r *mux.Router) {

	// user routes
	r.Handle("/api/user/create", GetCreateUserHandler()).Methods("POST")
	r.Handle("/api/user/login", GetLoginUserHandler()).Methods("POST")
	r.Handle("/api/user/delete", GetDeleteUserHandler()).Methods("DELETE")
	r.Handle("/api/user/update", GetUpdateUserHandler()).Methods("POST")
	r.Handle("/api/user/getAll", GetAllUsersHandler()).Methods("GET")
	r.Handle("/api/user/find/{id}", GetReadUserHandler()).Methods("GET")
	r.Handle("/api/user/like/{id}", GetLikeProfileHandler()).Methods("POST")

	// matches routes
	r.Handle("/api/match/delete/{mid}/{bid}", GetDeleteMatchHandler()).Methods("POST")

	// message routes
	r.Handle("/api/message/add", GetAddMessageHandler()).Methods("POST")
	r.Handle("/api/auth/zoom/callback", GetZoomHandler()).Methods("POST")
	r.Handle("/api/auth/zoom/tokencallback", GetZoomTokenHandler()).Methods("POST")
}

func CreatePostRoutes(r *mux.Router) {

	//r.HandleFunc("/api/user/test", user.PostTest).Methods("POST")
}
