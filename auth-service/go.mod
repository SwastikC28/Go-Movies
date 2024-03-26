module auth-service

go 1.22.0

replace shared => ../shared

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.1
	github.com/jinzhu/gorm v1.9.16
	github.com/satori/go.uuid v1.2.0
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	golang.org/x/crypto v0.0.0-20191205180655-e7c4368fe9dd // indirect
)
