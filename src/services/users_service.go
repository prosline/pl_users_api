package services

import (
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	return &user, nil
}

func GetUser(id int64) (*users.User, *errors.RestErr){
	return &users.User{}, nil
}
