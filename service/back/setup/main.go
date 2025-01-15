package main

import (
	"fmt"
	"os"
	"path"

	"github.com/eavlongs/file_sync/config"
	"github.com/eavlongs/file_sync/constants"
	"github.com/eavlongs/file_sync/models"
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
		db      = config.ConnectDatabase()
		_config = constants.NewConfig()
	)

	defer func() {
		config.CloseDatabaseConnection(db)
	}()

	var departments []models.Department

	db.Find(&departments)

	for _, department := range departments {
		err := os.MkdirAll(path.Join(_config.FileStoragePrefix, department.Name), 0755)

		if err != nil {
			fmt.Println(err)
		}
	}
}
