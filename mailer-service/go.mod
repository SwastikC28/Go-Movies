module mailer-service

go 1.22.0

replace shared => ../shared

require (
	github.com/go-sql-driver/mysql v1.8.0
	github.com/jinzhu/gorm v1.9.16
	github.com/vanng822/go-premailer v1.20.2
	github.com/xhit/go-simple-mail/v2 v2.16.0
)

require (
	github.com/PuerkitoBio/goquery v1.5.1 // indirect
	github.com/andybalholm/cascadia v1.1.0 // indirect
	github.com/go-test/deep v1.1.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/toorop/go-dkim v0.0.0-20201103131630-e1cd1a0a5208 // indirect
	github.com/vanng822/css v1.0.1 // indirect
	golang.org/x/net v0.21.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	shared v0.0.0-00010101000000-000000000000
)
