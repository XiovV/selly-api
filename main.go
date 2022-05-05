package main

import (
	"github.com/XiovV/selly-api/server"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

const (
	version = "0.1.0"
)

func main() {
	checkEnvVars()

	var logger *zap.Logger
	if os.Getenv("ENV") == "DEV" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	srv := server.New(sugar)

	sugar.Infow("running", "port", os.Getenv("PORT"), "environment", os.Getenv("ENV"), "version", version)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), srv.Serve()))
}
