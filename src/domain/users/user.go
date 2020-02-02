package users

import (
	"github.com/prosline/pl_users_api/src/utils/errors"
	"strings"
)

type User struct {
	Id          int64  `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	DateCreated string `json:"date_created" db:"date_created"`
	Status      string `json:"status" db:"status"`
	Password    string `json:"password" db:"password"`
}

func (user *User) IsValid() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	if user.FirstName == "" {
		return errors.NewBadRequestError("invalid user first name")
	}
	if user.LastName == "" {
		return errors.NewBadRequestError("invalid user last name")
	}
	if user.Status == "" {
		return errors.NewBadRequestError("invalid Status")
	}
	if user.Password == "" {
		return errors.NewBadRequestError("invalid password")
	}

	return nil
}
