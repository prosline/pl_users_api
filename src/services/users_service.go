package services

import (
	"errors"
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_util/utils/rest_errors"
)

var (
	UserService userServiceInterface = &userService{}
)

type userService struct {
}

type userServiceInterface interface {
	GetUser(int64) (*users.User, *rest_errors.RestErr)
	CreateUser(user users.User) (*users.User, *rest_errors.RestErr)
	UpdateUser(isPartial bool, user users.User) (*users.User, *rest_errors.RestErr)
	DeleteUser(id int64) *rest_errors.RestErr
	SearchUsers(string) (users.Users, *rest_errors.RestErr)
	Login(users.UserLogin) (*users.User, *rest_errors.RestErr)
}

func (s *userService) GetUser(id int64) (*users.User, *rest_errors.RestErr) {
	if id <= 0 {
		return nil, rest_errors.UserNotFound("Invalid User Id!")
	}
	user := users.User{Id: id}
	if err := user.Get(); err != nil {
		return nil, rest_errors.UserNotFound("User Not Found!")
	}

	return &user, nil
}
func (s *userService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	if err := user.IsValid(); err != nil {
		return nil, err
	}
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) UpdateUser(isPartial bool, user users.User) (*users.User, *rest_errors.RestErr) {
	u, err := s.GetUser(user.Id)
	if err != nil {
		return nil, rest_errors.UserNotFound("User Not Found!")
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
		return nil, rest_errors.NewInternalServerError("User update failed.",errors.New("Update Failed"))
	}
	return u, nil
}
func (s *userService) DeleteUser(id int64) *rest_errors.RestErr {
	var user users.User
	user.Id = id
	return user.Delete()
}
func (s *userService) SearchUsers(param string) (users.Users, *rest_errors.RestErr) {
	u := &users.User{}
	return u.FindUsersByStatus(param)
}

func (s *userService) Login(login users.UserLogin) (*users.User, *rest_errors.RestErr) {
	u := &users.User{
		Email:    login.Email,
		Password: login.Password,
	}
	if err := u.FindByEmailAndPassword(); err != nil {
		return nil, rest_errors.NewNotFoundError(err.Message)
	}
	return u, nil

}
