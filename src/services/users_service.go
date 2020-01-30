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
