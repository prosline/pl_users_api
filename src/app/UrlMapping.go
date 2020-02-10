package app

import (
	"github.com/prosline/pl_users_api/src/controllers"
)

func URLMapping() {
	router.GET("/ping", controllers.Ping)

	router.POST("/users", controllers.Create)
	router.GET("/users/:user_id", controllers.Get)
	router.PUT("/users/:user_id", controllers.Update)
	router.PATCH("/users/:user_id", controllers.Update)
	router.DELETE("/users/:user_id", controllers.Delete)
	router.GET("/internal/users/search", controllers.Search)
	router.POST("/users/login", controllers.Login)
}
