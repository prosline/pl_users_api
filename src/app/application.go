package app

import (
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_logger/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	logger.Info("Starting Application....")
	URLMapping()
	router.Run(":8080")
}
