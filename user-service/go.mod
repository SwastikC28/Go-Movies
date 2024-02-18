module user-service

go 1.22.0

replace shared => ../shared

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.1
	github.com/jinzhu/gorm v1.9.16
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	shared v0.0.0-00010101000000-000000000000
)
