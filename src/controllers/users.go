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
	UserNotCreated  = "Problem occurred while creating User record!"
	InvalidUserId   = "Invalid User Id"
	InvalidJSONBody = "Invalid Json Body"
	UserNotUpdated  = "Unable to update user"
)

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError(InvalidUserId)
		c.JSON(err.Status, err)
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

func CreateUser(c *gin.Context) {
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

func UpdateUser(c *gin.Context) {

	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError(InvalidUserId)
		c.JSON(err.Status, err)
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		er := errors.NewBadRequestError(InvalidJSONBody)
		c.JSON(er.Status, er)
	}
	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch
	if isPartial {

	}
	u, updateErr := services.UpdateUser(user)
	if updateErr != nil {
		er := errors.NewBadRequestError(UserNotUpdated)
		c.JSON(er.Status, er)
	}
	c.JSON(http.StatusOK, u)

}
