package users

import (
	"encoding/json"
	"fmt"
	"github.com/prosline/pl_logger/logger"
	"github.com/prosline/pl_util/utils/crypto"
	"github.com/prosline/pl_util/utils/rest_errors"
	"strings"
)

type Users []User

type PublicUser struct {
	Id          int64  `json:"id" db:"id"`
	DateCreated string `json:"date_created" db:"date_created"`
	Status      string `json:"status" db:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	DateCreated string `json:"date_created" db:"date_created"`
	Status      string `json:"status" db:"status"`
}

type User struct {
	Id          int64  `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	DateCreated string `json:"date_created" db:"date_created"`
	Status      string `json:"status" db:"status"`
	Password    string `json:"password" db:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) IsValid() *rest_errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("invalid email address")
	}
	if user.FirstName == "" {
		return rest_errors.NewBadRequestError("invalid user first name")
	}
	if user.LastName == "" {
		return rest_errors.NewBadRequestError("invalid user last name")
	}
	if user.Status == "" {
		return rest_errors.NewBadRequestError("invalid Status")
	}
	if user.Password == "" {
		return rest_errors.NewBadRequestError("invalid password")
	}

	return nil
}
func (users Users) Marshall(isPublic bool) []interface{} {
	rUser := make([]interface{}, len(users))
	for i, v := range users {
		rUser[i] = v.Marshall(isPublic)
	}
	return rUser
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}
	j, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(j, &privateUser)
	return privateUser
}

//func (uLogin *User) Marshall() interface{} {
//	return uLogin.Marshall(false)
//}

func (user *User) GetHashedPassword(pwrd string) string {
	hashedPassword, errCrypt := crypto.HashPassword(pwrd)
	if errCrypt != nil {
		logger.Info(fmt.Sprintf("Problem occurred while hashing password: $v \n", errCrypt.Error()))
		return ""
	}
	return strings.TrimSpace(hashedPassword)
}
