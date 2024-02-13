package main

import (
	"log"
	"time"
	"user-service/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db := connectDB()

	if db == nil {
		log.Println("Failed to Connect to DB")
	}

	app := config.NewApp("user-service", db, nil)

	app.Init()

	log.Println("Server started successfully.")
	app.StartServer()
}

func connectDB() *gorm.DB {
	url := "root:admin@tcp(localhost:3306)/go-movies?charset=utf8mb4&parseTime=true"

	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Println(err)
		return nil
	}

	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(500)
	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	db.LogMode(true)
	// utf8_general_ci is the default collate for utf8 and it is okay to not specify it.
	// ci means case insensitive.
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")

	db.BlockGlobalUpdate(true)

	log.Println("Connected to DB Successfully.")
	return db
}
