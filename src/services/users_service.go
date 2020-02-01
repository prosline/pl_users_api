package services

import (
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/utils/errors"
)

func GetUser(id int64) (*users.User, *errors.RestErr) {
	if id <= 0 {
		return nil, errors.UserNotFound("Invalid User Id!")
	}
	user := users.User{Id: id}
	if err := user.Get(); err != nil {
		return nil, errors.UserNotFound("User Not Found!")
	}

	return &user, nil
}
func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.IsValid(); err != nil {
		return nil, err
	}
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	u, err := GetUser(user.Id)
	if err != nil {
		return nil, errors.UserNotFound("User Not Found!")
	}
	if isPartial { // Handles the PATCH verb
		if user.FirstName != "" {
			u.FirstName = user.FirstName
		}
		if user.LastName != "" {
			u.LastName = user.LastName
		}
		if user.Email != "" {
			u.Email = user.Email
		}
	} else { // Handles the PUP verb
		u.Id = user.Id
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Email = user.Email
	}

	if er := u.Update(); er != nil {
		return nil, errors.NewInternalServerError("User update failed.")
	}
	return u, nil
}
func DeleteUser(id int64) *errors.RestErr {
	var user users.User
	user.Id = id
	return user.Delete()
}
