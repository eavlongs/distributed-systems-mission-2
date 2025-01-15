package routes

import (
	"github.com/eavlongs/file_sync/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterMainRoutes(router *gin.RouterGroup, c *controllers.MainController) {
	router.GET("files/*id", c.ServeFile)
}
