package main

import (
	"context"
	"fmt"
	"github.com/MaximPolyaev/gofermart/internal/adapters/accrual"
	"log"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/adapters/router"
	"github.com/MaximPolyaev/gofermart/internal/adapters/storage"
	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
	"github.com/MaximPolyaev/gofermart/internal/usecases/authusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/balanceusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/ordersusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/userusecase"
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

	accrualAdapter := accrual.New(*cfg.AccrualSystemAddress)

	ordersUseCase := ordersusecase.New(store, accrualAdapter)

	rtr := router.New(
		router.WithAuthUseCase(authusecase.New(store)),
		router.WithOrdersUseCase(ordersUseCase),
		router.WithUserUseCase(userusecase.New(store)),
		router.WithBalanceUseCase(balanceusecase.New(store)),
	)

	rtr.Configure()

	err = ordersUseCase.StartUpdateOrdersProcess(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("start server on", *cfg.RunAddress)
	if err := http.ListenAndServe(*cfg.RunAddress, rtr); err != nil {
		log.Fatal(err)
	}
}
