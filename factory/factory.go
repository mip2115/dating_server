package factory

import (
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
)

//NewUserVideoItem -
func NewUserVideoItem() *types.UserVideoItem {
	userUUID, _ := uuid.NewV4()
	uuid, _ := uuid.NewV4()
	return &types.UserVideoItem{
		UserUUID: mapping.StrToPtr(userUUID.String()),
		UUID:     mapping.StrToPtr(uuid.String()),
	}

}

// NewImage -
func NewImage() *types.Image {
	imageUUID, _ := uuid.NewV4()
	userUUID, _ := uuid.NewV4()
	link := "some-link"
	key := "some-key"
	return &types.Image{
		UUID:     mapping.StrToPtr(imageUUID.String()),
		UserUUID: mapping.StrToPtr(userUUID.String()),
		Link:     &link,
		Key:      &key,
	}
}

func NewUser() *types.User {
	userUUID, _ := uuid.NewV4()
	uuid := userUUID.String()
	email := "someemail@email.com"
	firstName := "first"
	lastName := "last"
	mobile := "mobile"
	dob := "dob"
	gender := "male"
	age := int64(21)
	partnerGender := "female"
	city := "New York City"
	maximumAge := int64(30)
	minimumAge := int64(20)
	zipcode := "10025"
	password := "password"

	return &types.User{
		UUID:          &uuid,
		Email:         &email,
		FirstName:     &firstName,
		LastName:      &lastName,
		Mobile:        &mobile,
		DOB:           &dob,
		Gender:        &gender,
		Age:           &age,
		PartnerGender: &partnerGender,
		City:          &city,
		MaximumAge:    &maximumAge,
		MinimumAge:    &minimumAge,
		Zipcode:       &zipcode,
		Password:      &password,
	}

}
