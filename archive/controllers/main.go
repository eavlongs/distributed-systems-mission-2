package controllers

import (
	"os"

	"github.com/eavlongs/file_sync/utils"
	"github.com/gin-gonic/gin"
)

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

func (c *MainController) ServeFile(ctx *gin.Context) {
	path := ctx.Param("id")

	// Remove the leading "/" from the wildcard parameter
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	ctx.File(path)
}

func (c *MainController) GetFiles(ctx *gin.Context) {
	files, err := os.ReadDir(".")
	if err != nil {
		utils.RespondWithInternalServerError(ctx, err.Error())
	}

	var filesToReturn []struct {
		Path string `json:"path"`
	}

	for _, file := range files {
		filesToReturn = append(filesToReturn, struct {
			Path string `json:"path"`
		}{Path: file.Name()})
	}

	utils.RespondWithSuccess(ctx, gin.H(map[string]interface{}{
		"files": filesToReturn,
	}))
}
