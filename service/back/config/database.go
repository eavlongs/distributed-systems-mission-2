package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//goland:noinspection ALL
func ConnectDatabase() *gorm.DB {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_DATABASE := os.Getenv("DB_DATABASE")
	DB_PORT := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Phnom_Penh", DB_HOST, DB_USERNAME, DB_PASSWORD, DB_DATABASE, DB_PORT)

	db, errDb := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		TranslateError: true,
	})

	if errDb != nil {
		panic("failed to connect database")
	}

	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	postgresDb, err := db.DB()

	if err != nil {
		println(err)
		panic(err)
	}

	postgresDb.Close()
}
