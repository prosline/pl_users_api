package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_oauth/oauth"
	"github.com/prosline/pl_users_api/src/domain/users"
	"github.com/prosline/pl_users_api/src/services"
	"github.com/prosline/pl_util/utils/rest_errors"
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

func getUserId(param string) (int64, rest_errors.RestErr) {
	userId, userErr := strconv.ParseInt(param, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError(InvalidUserId)
	}
	return userId, nil
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil{
		c.JSON(err.Status(), err)
		return
	}
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	user, getUserErr := services.UserService.GetUser(userId)
	if getUserErr != nil {
		err := rest_errors.CreateUserError(UserNotCreated)
		c.JSON(err.Status(), getUserErr)
		return
	}
	if oauth.GetUserId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPlublic(c.Request)))
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "Implementation required!")
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError(InvalidJSONBody)
		c.JSON(restErr.Status(), restErr)
		return
	}
	u, errUser := services.UserService.CreateUser(user)
	if errUser != nil {
		c.JSON(http.StatusBadRequest, errUser)
		return
	}
	c.JSON(http.StatusOK, u.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		er := rest_errors.NewBadRequestError(InvalidJSONBody)
		c.JSON(er.Status(), er)
		return
	}
	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch
	u, updateErr := services.UserService.UpdateUser(isPartial, user)
	if updateErr != nil {
		er := rest_errors.NewBadRequestError(UserNotUpdated)
		c.JSON(er.Status(), er)
		return
	}
	c.JSON(http.StatusOK, u.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	userId, err := getUserId(c.Param("user_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	delErr := services.UserService.DeleteUser(userId)
	if delErr != nil {
		delErr := rest_errors.NewBadRequestError(UserNotDeleted)
		c.JSON(delErr.Status(), delErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

// '2016-06-22 19:10:25-07'
func Search(c *gin.Context) {
	param := c.Query("status")

	if param == "" {
		paramErr := rest_errors.NewBadRequestError(UserNotDeleted)
		c.JSON(paramErr.Status(), paramErr)
	}
	allUsers, srchErr := services.UserService.SearchUsers(param)
	if srchErr != nil {
		err := rest_errors.UserNotFound(UserNotAvailable)
		c.JSON(err.Status(), err)
	}
	c.JSON(http.StatusOK, allUsers.Marshall(c.GetHeader("X-Public") == "true"))
}
func Login(c *gin.Context) {
	var userLogin users.UserLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		er := rest_errors.NewBadRequestError("Invalid JSON request during Login process")
		c.JSON(http.StatusBadRequest, er)
		return
	}
	login, err := services.UserService.Login(userLogin)
	if err != nil {
		er := rest_errors.UserNotFound("Incorrect Login and Password information")
		//c.JSON(http.StatusBadRequest, "Incorrect Login and Password information")
		c.JSON(http.StatusBadRequest, er)
		return
	}
	c.JSON(http.StatusOK, login.Marshall(c.GetHeader("X-Public") == "true"))
}
