package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eavlongs/file_sync/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MainMiddleware struct {
	db *gorm.DB
}

func NewMainMiddleware(db *gorm.DB) *MainMiddleware {
	return &MainMiddleware{db: db}
}

func (m *MainMiddleware) TestFail(ctx *gin.Context) {
	utils.RespondWithUnauthorizedError(ctx)
	ctx.Abort()
}

func (m *MainMiddleware) TestSuccess(ctx *gin.Context) {
	ctx.Next()
}

func (m *MainMiddleware) IsUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := getLoggedInUserClaims(ctx)

		if err != nil {
			utils.RespondWithBadRequestError(ctx, err.Error())
			ctx.Abort()
		}
		ctx.Set("_auth_user_id", claims.ID)
		ctx.Set("_auth_user_department_id", claims.DepartmentID)
		ctx.Next()
	}
}

func (m *MainMiddleware) WhoAmI(ctx *gin.Context) {

	claims, err := getLoggedInUserClaims(ctx)

	if err != nil {
		utils.RespondWithBadRequestError(ctx, err.Error())
		return
	}

	utils.RespondWithSuccess(ctx, claims)
}

func getLoggedInUserClaims(ctx *gin.Context) (*utils.Claims, error) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.New("authorization header missing or invalid")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
