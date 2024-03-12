module movie-service

go 1.22.0

replace shared => ../shared

require (
	github.com/go-sql-driver/mysql v1.8.0
	github.com/gorilla/mux v1.8.1
	github.com/jinzhu/gorm v1.9.16
	github.com/satori/go.uuid v1.2.0
	shared v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cloudinary/cloudinary-go v1.7.0 // indirect
	github.com/cloudinary/cloudinary-go/v2 v2.7.0 // indirect
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/schema v1.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	golang.org/x/crypto v0.21.0 // indirect
)
