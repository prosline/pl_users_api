package app

import (
	"github.com/prosline/pl_users_api/src/controllers"
)

func URLMapping() {
	router.GET("/ping", controllers.Ping)
	router.GET("/users/:user_id", controllers.GetUser)
	//	router.GET("/users/search/:id", controllers.FindUser)
	router.POST("/users", controllers.CreateUser)

}
