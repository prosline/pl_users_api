package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/services"
	"github.com/prosline/pl_users_api/src/utils/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	BAD_REQUEST       = "Invalid JSON body, please review user info!"
	BAD_EMAIL_REQUEST = "Invalid Email Address"
	USER_NOT_CREATED  = "Problem occurred while creating a User record!"
)

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user Id!")
		c.JSON(err.Status, err)
	}
	user, getUserErr := services.GetUser(userId)
	if getUserErr != nil {
		err := errors.CreateUserError(USER_NOT_CREATED)
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
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	if er := json.Unmarshal(data, &user); er != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	//	if err := c.ShouldBindJSON(&user); err != nil {
	//		restErr := errors.NewBadRequestError("invalid json body")
	//		c.JSON(restErr.Status, restErr)
	//		return
	//	}
	result, errUser := services.CreateUser(user)
	if errUser != nil {
		c.JSON(http.StatusBadRequest, errUser)
		return
	}
	c.JSON(http.StatusCreated, result)
}
