package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eavlongs/file_sync/constants"
	"github.com/eavlongs/file_sync/models"
	"github.com/eavlongs/file_sync/repository"
	"github.com/eavlongs/file_sync/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MainController struct {
	repo *repository.MainRepository
	cf   *constants.Config
}

func NewMainController(repo *repository.MainRepository, st *constants.Config) *MainController {
	return &MainController{repo: repo, cf: st}
}

func (c *MainController) Test(ctx *gin.Context) {
	utils.RespondWithSuccess(ctx, struct {
		Data string `json:"data"`
	}{Data: "Hello World"})
}

func (c *MainController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.RespondWithBadRequestError(ctx, "File not found in body")
		return
	}

	// generate unique id
	var id string

	for {
		id = uuid.New().String()

		_, err := c.repo.GetFileById(id)

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			break
		} else {
			fmt.Println("Error checking if id exists in database:", err)
		}
	}

	fileExtension := filepath.Ext(file.Filename)
	fileName := id + fileExtension

	departmentID := ctx.Keys["_auth_user_department_id"].(uint)
	department, err := c.repo.GetDepartmentByID(departmentID)

	if err != nil {
		fmt.Println("Error getting department from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get department from database")
		return
	}

	fi, err := file.Open()

	if err != nil {
		fmt.Println("Error opening file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to open file")
		return
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	filePath := fmt.Sprintf("%s/%s/%s", c.cf.FileStoragePrefix, department.Name, fileName)
	fo, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to create file")
		return
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(fo, fi)
	if err != nil {
		fmt.Println("Error copying file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to copy file")
		return
	}

	// send file to archive server
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to create form file")
		return
	}

	fi.Seek(0, 0)

	_, err = io.Copy(part, fi)
	if err != nil {
		fmt.Println("Error copying file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to copy file")
		return
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to close writer")
		return
	}

	req, err := http.NewRequest("POST", "http://"+c.cf.ServerMapping[constants.ARCHIVE_SERVER]+"api/upload", &requestBody)

	if err != nil {
		fmt.Println("Error creating request:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending replication request")
	} else if resp.StatusCode != 200 {
		// check response status
		fmt.Println("Error sending file to archive server")
	}

	// create entry in database

	err = c.repo.CreateFile(id, file.Filename, filePath, departmentID)

	if err != nil {
		fmt.Println("Error inserting file entry into database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to insert file entry into database")
		return
	}

	utils.RespondWithSuccess(ctx, nil)
}

func (c *MainController) EditFile(ctx *gin.Context) {
	fileID := ctx.Param("id")
	dbFile, err := c.repo.GetFileById(fileID)

	if err != nil {
		fmt.Println("Error getting file from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get file from database")
		return
	}

	departmentID := ctx.Keys["_auth_user_department_id"].(uint)

	if dbFile.DepartmentID != departmentID {
		utils.RespondWithUnauthorizedError(ctx)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		utils.RespondWithBadRequestError(ctx, "File not found in body")
		return
	}

	fi, err := file.Open()

	if err != nil {
		fmt.Println("Error opening file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to open file")
		return
	}

	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	filePath := dbFile.Path

	fo, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to create file")
		return
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(fo, fi)
	if err != nil {
		fmt.Println("Error copying file:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to copy file")
		return
	}

	err = c.repo.EditFile(fileID)

	if err != nil {
		fmt.Println("Error inserting file entry into database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to insert file entry into database")
		return
	}

	utils.RespondWithSuccess(ctx, nil)
}

func (c *MainController) DeleteFile(ctx *gin.Context) {
	fileID := ctx.Param("id")
	dbFile, err := c.repo.GetFileById(fileID)

	if err != nil {
		fmt.Println("Error getting file from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get file from database")
		return
	}

	departmentID := ctx.Keys["_auth_user_department_id"].(uint)

	if dbFile.DepartmentID != departmentID {
		utils.RespondWithUnauthorizedError(ctx)
		return
	}

	err = os.Remove(dbFile.Path)

	if err != nil {
		fmt.Println("Error deleting file from disk:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to delete file from disk")
		return
	}

	err = c.repo.DeleteFile(fileID)

	if err != nil {
		fmt.Println("Error deleting file from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to delete file from database")
		return
	}

	utils.RespondWithSuccess(ctx, nil)
}

func (c *MainController) getFiles(departmentID uint, serverType int) ([]models.File, error) {
	files, err := c.repo.GetFiles(departmentID, serverType)

	if err != nil {
		return nil, err
	}

	for i := range files {
		files[i].Path = c.cf.ServerMapping[serverType] + files[i].ID
	}

	return files, nil
}

func (c *MainController) GetFiles(ctx *gin.Context) {
	departmentID := ctx.Keys["_auth_user_department_id"].(uint)

	files, err := c.getFiles(departmentID, constants.MAIN_SERVER)

	if err != nil {
		fmt.Println("Error getting files from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get files from database")
		return
	}

	utils.RespondWithSuccess(ctx, struct {
		Files []models.File `json:"files"`
	}{Files: files})
}

func (c *MainController) GetBackupFiles(ctx *gin.Context) {
	departmentID := ctx.Keys["_auth_user_department_id"].(uint)

	files, err := c.getFiles(departmentID, constants.BACKUP_SERVER)

	if err != nil {
		fmt.Println("Error getting files from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get files from database")
		return
	}

	utils.RespondWithSuccess(ctx, struct {
		Files []models.File `json:"files"`
	}{Files: files})
}

func (c *MainController) GetArchiveFiles(ctx *gin.Context) {
	departmentID := ctx.Keys["_auth_user_department_id"].(uint)

	files, err := c.getFiles(departmentID, constants.ARCHIVE_SERVER)

	if err != nil {
		fmt.Println("Error getting files from database:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to get files from database")
		return
	}

	utils.RespondWithSuccess(ctx, struct {
		Files []models.File `json:"files"`
	}{Files: files})
}

func (c *MainController) SyncNow(ctx *gin.Context) {
	err := exec.Command(c.cf.SyncFileScriptPath).Run()

	if err != nil {
		fmt.Println("Error syncing files:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to sync files")
		return
	}

	utils.RespondWithSuccess(ctx, nil)
}

func (c *MainController) GetDepartments(ctx *gin.Context) {
	departments, err := c.repo.GetDepartments()

	if err != nil {
		utils.RespondWithInternalServerError(ctx, "Couldn't get departments")
		return
	}

	utils.RespondWithSuccess(ctx, departments)
}

func (c *MainController) ServeFile(ctx *gin.Context) {
	id := ctx.Param("id")
	// check if id exists in database

	file, err := c.repo.GetFileById(id)

	if err != nil {
		utils.RespondWithBadRequestError(ctx, err.Error())
	}

	ctx.File(file.Path)
}

func (c *MainController) SyncFile(ctx *gin.Context) {
	err := exec.Command(c.cf.SyncFileScriptPath).Run()

	if err != nil {
		fmt.Println("Error syncing files:", err)
		utils.RespondWithInternalServerError(ctx, "Failed to sync files")
		return
	}

	utils.RespondWithSuccess(ctx, nil)
}
