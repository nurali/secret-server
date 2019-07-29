package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/nurali/secret-server/secret-service/pkg/app"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	"github.com/nurali/secret-server/secret-service/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()

	initLogger(cfg.GetLogLevel())

	// setup DB
	var db *gorm.DB
	var err error
	if db, err = app.OpenDB(cfg); err != nil {
		log.Panicf("Issue with connecting to DB, %v", err)
	}
	if db != nil {
		err = db.DB().Ping()
		if err != nil {
			log.Panicf("Unable to connect to DB, %v", err)
		}
		defer db.Close()
	}
	err = app.SetupDB(db)
	if err != nil {
		log.Panicf("Unable to setup DB, %v", err)
	}
	log.Infof("Database OK")

	// setup app
	router := app.Router(db)
	startService(cfg.GetHttpPort(), router)
}

func initLogger(logLevel string) {
	level, _ := log.ParseLevel(logLevel)
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
}

func startService(port int, router *mux.Router) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Infof("secret service running at:%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start secret service, error:%v", err)
	}
}
