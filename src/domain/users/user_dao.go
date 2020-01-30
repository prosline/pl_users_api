package users

import (
	"fmt"
	"github.com/prosline/pl_users_api/src/utils/errors"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Save() *errors.RestErr {
	usr := userDB[user.Id]
	if usr != nil {
		err := errors.NewUserBadRequest(fmt.Sprintf("Bad Request, user %d / %s already exists", user.Id, user.FirstName))
		return err
	}
	userDB[user.Id] = user
	return nil
}

func (user *User) Get() *errors.RestErr {
	usr := userDB[user.Id]
	if usr == nil {
		err := errors.UserNotFound(fmt.Sprintf("User id %d not found", user.Id))
		return err
	}
	user.Id = usr.Id
	user.FirstName = usr.FirstName
	user.LastName = usr.LastName
	user.Email = usr.Email

	return nil
}
