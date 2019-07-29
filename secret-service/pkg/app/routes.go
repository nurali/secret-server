package app

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/nurali/secret-server/secret-service/pkg/ctrl"
)

func Router(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	secretCtrl := ctrl.NewSecretCtrl(db)
	// TODO use secrets instead of secret
	r.HandleFunc("/api/secret", secretCtrl.Create).Methods("POST")
	r.HandleFunc("/api/secret/{hash}", secretCtrl.Show).Methods("GET")

	return r
}
