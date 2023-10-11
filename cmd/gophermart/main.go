package main

import (
	"log"

	"github.com/MaximPolyaev/gofermart/internal/adapters/http"
	"github.com/MaximPolyaev/gofermart/internal/adapters/router"
	"github.com/MaximPolyaev/gofermart/internal/adapters/storage"
	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
	"github.com/MaximPolyaev/gofermart/internal/migration"
	"github.com/MaximPolyaev/gofermart/internal/usecases/authusecase"
)

func main() {
	cfg := config.New()
	if err := cfg.Parse(); err != nil {
		log.Fatal(err)
	}

	db, err := dbconn.InitDB(*cfg.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.MigrateUp(db); err != nil {
		log.Fatal(err)
	}

	server := http.New(*cfg.RunAddress)

	store := storage.New(db)

	auth := authusecase.New(store)

	rtr := router.New(auth)
	rtr.Configure()

	if err := server.ListenAndServe(rtr.Mux); err != nil {
		log.Fatal(err)
	}
}
