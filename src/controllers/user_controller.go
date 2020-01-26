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

func GetUser(c *gin.Context) {
	//TODO: Get parameter from request object
	id, _ := strconv.ParseInt(c.Param("user_id"),10,64)

	//TODO: Request user from user service
	result, err := services.GetUser(id)
	//TODO: Send user using response object
	fmt.Printf("Id = %d Err = %v", result, err)

}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implementation required!")

}

func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		err := errors.BadRequestError("Invalid JSON error")
		c.JSON(err.Code,err)
		return
	}
	result, userErr := services.CreateUser(user)
	if userErr != nil {
		err := errors.BadRequestError("Create User Error")
		c.JSON(err.Code,err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
