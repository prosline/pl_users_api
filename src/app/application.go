package app

import (
	"github.com/gin-gonic/gin"
	"github.com/prosline/pl_logger/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	logger.Info("Starting Application on Port 8081....")
	URLMapping()
	router.Run(":8081")
}
