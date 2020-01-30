package users

import (
	"github.com/prosline/pl_users_api/src/utils/errors"
	"strings"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
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

	return nil
}
