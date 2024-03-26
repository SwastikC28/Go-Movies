package main

import (
	"log"
	"mailer-service/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var counts int64

func main() {
	var db *gorm.DB = nil
	db = connectDB()
	defer db.Close()

	if db == nil {
		log.Println("Failed to Connect to DB")
		panic("Failed to connect to DB")
	}

	mailerMS := config.NewApp("mailer-service", db, &sync.WaitGroup{})

	// Initialize User Microservice
	mailerMS.Init()

	// Start Microservice
	log.Println("Mailer microservice started successfully.")
	go mailerMS.StartServer()

	config.StartEventHandler()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	mailerMS.StopServer()
	config.StopEventHandler()
}

func connectDB() *gorm.DB {
	// url := "root:admin@tcp(localhost:3306)/go-movies?charset=utf8mb4&parseTime=true"
	url := "admin:admin@tcp(mariadb)/go-movies?charset=utf8mb4&parseTime=True&loc=Local"

	for {
		db, err := gorm.Open("mysql", url)

		if err != nil {
			log.Println("SQLDB not yet ready ...")
			counts++
		} else {
			sqlDB := db.DB()
			sqlDB.SetMaxIdleConns(500)
			sqlDB.SetMaxOpenConns(2)
			sqlDB.SetConnMaxLifetime(3 * time.Minute)

			db.LogMode(true)
			// utf8_general_ci is the default collate for utf8 and it is okay to not specify it.
			// ci means case insensitive.
			// db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")

			db.BlockGlobalUpdate(true)

			log.Println("Connected to DB successfully.")
			return db
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}

}
