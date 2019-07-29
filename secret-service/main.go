package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nurali/secret-server/secret-service/pkg/app"

	"github.com/gorilla/mux"

	"github.com/nurali/secret-server/secret-service/pkg/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()

	initLogger(cfg.GetLogLevel())

	router := app.Router()

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
