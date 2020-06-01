package routes

import (
	"../auth"
	"../endpoints/user"
	"net/http"
)

func GetCreateUserHandler() http.Handler {
	handler := http.HandlerFunc(user.CreateUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetUpdateUserHandler() http.Handler {
	handler := http.HandlerFunc(user.UpdateUser)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}

func GetLoginUserHandler() http.Handler {
	handler := http.HandlerFunc(user.LoginUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetDeleteUserHandler() http.Handler {
	handler := http.HandlerFunc(user.DeleteUser)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}

func GetReadUserHandler() http.Handler {
	handler := http.HandlerFunc(user.GetUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetAllUsersHandler() http.Handler {
	handler := http.HandlerFunc(user.GetAllUsers)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetLikeProfileHandler() http.Handler {
	handler := http.HandlerFunc(user.LikeProfile)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}
