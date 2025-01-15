package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eavlongs/file_sync/config"
	"github.com/eavlongs/file_sync/constants"
	"github.com/eavlongs/file_sync/controllers"
	"github.com/eavlongs/file_sync/middlewares"
	"github.com/eavlongs/file_sync/repository"
	"github.com/eavlongs/file_sync/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	var (
		db     = config.ConnectDatabase()
		Config = constants.NewConfig()
		// db             *gorm.DB
		mainMiddleware = middlewares.NewMainMiddleware(db)
		AuthRepository = repository.NewAuthRepository(db)
		MainRepository = repository.NewMainRepository(db)

		AuthController = controllers.NewAuthController(AuthRepository)
		MainController = controllers.NewMainController(MainRepository, Config)
	)

	defer func() {
		config.CloseDatabaseConnection(db)
	}()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}

	router.Use(cors.New(corsConfig))

	routePrefix := os.Getenv("API_ROUTE_PREFIX")
	routerGroup := router.Group(routePrefix)

	routes.RegisterAuthRoutes(routerGroup, AuthController, mainMiddleware)
	routes.RegisterMainRoutes(routerGroup, MainController, mainMiddleware)

	port := os.Getenv("API_PORT")
	if err := router.Run("127.0.0.1:" + port); err != nil {
		// if err := route.Run(":" + port); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
