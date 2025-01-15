package controllers

import (
	"github.com/eavlongs/file_sync/models"
	"github.com/eavlongs/file_sync/repository"
	"github.com/eavlongs/file_sync/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo *repository.AuthRepository
}

func NewAuthController(repo *repository.AuthRepository) *AuthController {
	return &AuthController{repo: repo}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.RespondWithBadRequestError(ctx, "Invalid request")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithInternalServerError(ctx, "Error hashing password")
		return
	}
	user.Password = string(hashedPassword)

	if err := c.repo.FindDepartmentByID(user.DepartmentID); err != nil {
		utils.RespondWithBadRequestError(ctx, "Invalid department ID")
		return
	}

	if err := c.repo.CreateUser(&user); err != nil {
		utils.RespondWithInternalServerError(ctx, "Error creating user")
		return
	}

	token, err := utils.GenerateJWT(&user)
	if err != nil {
		utils.RespondWithInternalServerError(ctx, "Error generating token")
		return
	}

	utils.RespondWithSuccess(ctx, map[string]string{"token": token})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithBadRequestError(ctx, "Invalid request")
		return
	}

	var user models.User
	if err := c.repo.FindUser(input.Email, &user); err != nil {
		utils.RespondWithBadRequestError(ctx, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.RespondWithBadRequestError(ctx, "Invalid email or password")
		return
	}

	token, err := utils.GenerateJWT(&user)
	if err != nil {
		utils.RespondWithInternalServerError(ctx, "Error generating token")
		return
	}

	utils.RespondWithSuccess(ctx, map[string]string{"token": token})
}
