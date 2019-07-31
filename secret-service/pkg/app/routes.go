package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nurali/secret-server/secret-service/pkg/ctrl"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func Router(db *gorm.DB, middleware Middleware) *mux.Router {
	r := mux.NewRouter()

	secretCtrl := ctrl.NewSecretCtrl(db)
	// TODO use secrets instead of secret
	r.HandleFunc("/api/secret", middleware(secretCtrl.Create)).Methods("POST")
	r.HandleFunc("/api/secret/{hash}", middleware(secretCtrl.Show)).Methods("GET")

	return r
}
