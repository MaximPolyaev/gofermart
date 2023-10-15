package main

import (
	"fmt"
	"github.com/MaximPolyaev/gofermart/internal/usecases/userusercase"
	"log"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/adapters/router"
	"github.com/MaximPolyaev/gofermart/internal/adapters/storage"
	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
	"github.com/MaximPolyaev/gofermart/internal/usecases/authusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/ordersusecase"
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

	store := storage.New(db)

	rtr := router.New(
		router.WithAuthUseCase(authusecase.New(store)),
		router.WithOrdersUseCase(ordersusecase.New(store)),
		router.WithUserUseCase(userusercase.New(store)),
	)

	rtr.Configure()

	fmt.Println("start server on", *cfg.RunAddress)
	if err := http.ListenAndServe(*cfg.RunAddress, rtr); err != nil {
		log.Fatal(err)
	}
}
