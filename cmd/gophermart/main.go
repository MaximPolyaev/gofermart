package main

import (
	"log"

	"github.com/MaximPolyaev/gofermart/internal/adapters/http"
	"github.com/MaximPolyaev/gofermart/internal/adapters/router"
	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
)

func main() {
	cfg := config.New()
	if err := cfg.Parse(); err != nil {
		log.Fatal(err)
	}

	_, err := dbconn.InitDB(*cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	server := http.New(*cfg.RunAddress)

	rtr := router.New()
	rtr.Configure()

	if err := server.ListenAndServe(rtr.Mux); err != nil {
		log.Fatal(err)
	}
}
