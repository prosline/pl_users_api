package app

import (
	"github.com/prosline/pl_users_api/src/controllers"
)

func URLMapping() {
	router.GET("/ping", controllers.Ping)

	router.POST("/users", controllers.CreateUser)
	router.GET("/users/:user_id", controllers.GetUser)
	router.PUT("/users/:user_id", controllers.UpdateUser)
	router.PATCH("/users/:user_id", controllers.UpdateUser)

}
