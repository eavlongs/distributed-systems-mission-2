package routes

import (
	"github.com/eavlongs/file_sync/controllers"
	"github.com/eavlongs/file_sync/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterMainRoutes(router *gin.RouterGroup, c *controllers.MainController, m *middlewares.MainMiddleware) {
	router.GET("/", c.Test)

	router.POST("files/upload", m.IsUser(), c.UploadFile)
	router.GET("files/main", m.IsUser(), c.GetFiles)
	router.GET("files/backup", m.IsUser(), c.GetBackupFiles)
	router.GET("files/archive", m.IsUser(), c.GetArchiveFiles)
	router.GET("files/:id", c.ServeFile)

	router.DELETE("files/:id", m.IsUser(), c.DeleteFile)

	router.POST("files/sync", m.IsUser(), c.SyncFile)
	router.GET("/departments", c.GetDepartments)
}
