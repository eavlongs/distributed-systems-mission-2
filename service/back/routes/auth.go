package routes

import (
	"github.com/eavlongs/file_sync/controllers"
	"github.com/eavlongs/file_sync/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, c *controllers.AuthController, m *middlewares.MainMiddleware) {
	r := router.Group("/auth")
	r.POST("/register", c.Register)
	r.POST("/login", c.Login)

	r.GET("/whoami", m.WhoAmI)
}
