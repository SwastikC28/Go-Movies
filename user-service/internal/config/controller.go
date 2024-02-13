package config

import "github.com/gorilla/mux"

type Controller interface {
	RegisterRoute(router *mux.Router)
}
