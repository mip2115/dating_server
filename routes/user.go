package routes

import (
	//"code.mine/dating_server/auth"
	"net/http"

	"code.mine/dating_server/auth"
	"code.mine/dating_server/endpoints/user_handler"
)

func GetCreateUserHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.CreateUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetUpdateUserHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.UpdateUser)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}

func GetLoginUserHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.LoginUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetDeleteUserHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.DeleteUser)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}

func GetReadUserHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.GetUser)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetAllUsersHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.GetAllUsers)
	//h := auth.RefreshJWT(handler)
	//h = auth.VerifyJWT(h)
	return handler
}

func GetLikeProfileHandler() http.Handler {
	handler := http.HandlerFunc(user_handler.LikeProfile)
	h := auth.RefreshJWT(handler)
	h = auth.VerifyJWT(h)
	return h
}
