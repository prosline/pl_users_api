package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/services"
	"github.com/prosline/pl_users_api/src/utils/errors"
	"net/http"
	"strconv"
)

const (
	UserNotCreated   = "Problem occurred while creating User record!"
	InvalidUserId    = "Invalid User Id"
	InvalidJSONBody  = "Invalid Json Body"
	UserNotUpdated   = "Unable to update user"
	UserNotDeleted   = "Unable to delete user"
	UserNotAvailable = "Users not available"
)

func getUserId(param string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(param, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError(InvalidUserId)
	}
	return userId, nil
}

func Get(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	user, getUserErr := services.GetUser(userId)
	if getUserErr != nil {
		err := errors.CreateUserError(UserNotCreated)
		c.JSON(err.Status, getUserErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implementation required!")
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(InvalidJSONBody)
		c.JSON(restErr.Status, restErr)
		return
	}
	result, errUser := services.CreateUser(user)
	if errUser != nil {
		c.JSON(http.StatusBadRequest, errUser)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func Update(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		er := errors.NewBadRequestError(InvalidJSONBody)
		c.JSON(er.Status, er)
	}
	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch
	u, updateErr := services.UpdateUser(isPartial, user)
	if updateErr != nil {
		er := errors.NewBadRequestError(UserNotUpdated)
		c.JSON(er.Status, er)
	}
	c.JSON(http.StatusOK, u)
}

func Delete(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	delErr := services.DeleteUser(userId)
	if delErr != nil {
		delErr := errors.NewBadRequestError(UserNotDeleted)
		c.JSON(delErr.Status, delErr)
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// '2016-06-22 19:10:25-07'
func Search(c *gin.Context) {
	param := c.Query("status")

	if param == "" {
		paramErr := errors.NewBadRequestError(UserNotDeleted)
		c.JSON(paramErr.Status, paramErr)
	}
	allUsers, srchErr := services.SearchUsers(param)
	if srchErr != nil {
		err := errors.UserNotFound(UserNotAvailable)
		c.JSON(err.Status, err)
	}
	c.JSON(http.StatusOK, allUsers)

}
