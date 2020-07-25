package user_handler

import (
	"encoding/json"
	"net/http"

	"code.mine/dating_server/auth"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/service/aws_service"
	"code.mine/dating_server/service/image_service"
	"code.mine/dating_server/service/user_service"
	"code.mine/dating_server/types"
	"github.com/gorilla/mux"
)

func UploadUserImage(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	imageRequest := &types.ImageUploadRequest{}
	err := json.NewDecoder(r.Body).Decode(imageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imageFileBytes := imageRequest.FileBytes
	link, key, err := aws_service.UploadImageToS3(imageFileBytes, &userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	img, err := image_service.CreateImage(&userUUID, link, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = user_service.SaveUserImage(&userUUID, img.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.UploadImageResponse{}
	response.Token = &tkString
	response.UploadedImageLink = link
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUserImage(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	imageRequest := &types.ImageDeleteRequest{}
	err := json.NewDecoder(r.Body).Decode(imageRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = aws_service.DeleteImageFromS3(imageRequest.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = image_service.DeleteImage(imageRequest.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = user_service.RemoveUserImage(&userUUID, imageRequest.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.ImageDeleteResponse{}
	response.Token = &tkString
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TODO – log user in right away
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := user_service.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserID = res
	user.Password = nil
	response := types.UserResponse{}
	response.User = user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	returnedUser, err := user_service.LoginUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tkString, err := auth.GenerateJWT(returnedUser.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	response := types.UserResponse{}
	response.Token = &tkString
	response.User = returnedUser
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := user_service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//userUUID := r.Context().Value("userUUID").(string)
	tkString := r.Context().Value("tokenString").(string)
	user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = user_service.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.UserResponse{}
	response.Token = &tkString
	response.User = user
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LikeProfile(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userID").(string)
	tkString := r.Context().Value("tokenString").(string)
	params := mux.Vars(r)
	likedProfileUUID := params["likedProfileUUID"]
	// if m isnt' nil then a match occuredd
	m, err := user_service.LikeProfile(&userUUID, &likedProfileUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.MatchResponse{}
	response.Token = &tkString
	if m != nil {
		response.Match = m
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Context().Value("userUUID").(string)
	err := user_service.DeleteUser(&userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.UserResponse{}
	response.Op = mapping.BoolToPtr(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userUUID := params["uuid"]
	user, err := user_service.GetUser(&userUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := types.UserResponse{}
	response.User = user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
