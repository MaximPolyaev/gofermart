package main

import (
	"fmt"
	"log"
	"net/http"

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

	store := storage.New(db)

	auth := authusecase.New(store)

	rtr := router.New(auth)
	rtr.Configure()

	fmt.Println("start server on", *cfg.RunAddress)
	if err := http.ListenAndServe(*cfg.RunAddress, rtr); err != nil {
		log.Fatal(err)
	}
}
