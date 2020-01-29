package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/services"
	"github.com/prosline/pl_users_api/src/utils/errors"
	"net/http"
	"strconv"
)

const (
	BAD_REQUEST      = "Invalid JSON, please review user info!"
	USER_NOT_CREATED = "Problem occurred while creating a User record!"
)

func GetUser(c *gin.Context) {
	//TODO: Get parameter from request object
	//id, _ := strconv.ParseInt(c.Param("user_id"),10,64)
	id, _ := strconv.Atoi(c.Param("user_id"))

	//TODO: Request user from user service
	result, err := services.GetUser(int64(id))
	//TODO: Send user using response object
	fmt.Printf("Id = %d Err = %v \n", id, err)
	fmt.Printf("Result (user.email) = %v", result.Email)

}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implementation required!")
}

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errMsg := errors.BadRequestError(BAD_REQUEST)
		c.JSON(errMsg.Status, errMsg)
		return
	}
	result, userErr := services.CreateUser(user)
	if userErr != nil {
		errMsg := errors.CreateUserError(USER_NOT_CREATED)
		c.JSON(errMsg.Status, errMsg)
		return
	}
	c.JSON(http.StatusCreated, result)
}
