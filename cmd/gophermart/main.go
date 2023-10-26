package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaximPolyaev/gofermart/internal/adapters/accrual"
	"github.com/MaximPolyaev/gofermart/internal/adapters/router"
	"github.com/MaximPolyaev/gofermart/internal/adapters/storage"
	"github.com/MaximPolyaev/gofermart/internal/config"
	"github.com/MaximPolyaev/gofermart/internal/dbconn"
	"github.com/MaximPolyaev/gofermart/internal/usecases/accrualusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/authusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/balanceusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/ordersusecase"
	"github.com/MaximPolyaev/gofermart/internal/usecases/userusecase"
	"github.com/MaximPolyaev/gofermart/internal/utils/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.New(os.Stdout)

	err := run(context.Background(), log)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cfg := config.New()
	if err := cfg.Parse(); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"uri": *cfg.DatabaseURI,
	}).Info("init db connection")

	db, err := dbconn.InitDB(*cfg.DatabaseURI)
	if err != nil {
		return err
	}

	store := storage.New(db, log)

	log.WithFields(logrus.Fields{
		"address": *cfg.AccrualSystemAddress,
	}).Info("init accrual adapter")

	shutdownHandler(cancel, log)

	accrualAdapter := accrual.New(*cfg.AccrualSystemAddress, log)
	accrualUseCase := accrualusecase.New(accrualAdapter, store, log)

	go accrualUseCase.StartSyncOrdersStatusesProcess(ctx)

	rtr := router.New(
		router.WithAuthUseCase(authusecase.New(store)),
		router.WithOrdersUseCase(ordersusecase.New(store)),
		router.WithUserUseCase(userusecase.New(store)),
		router.WithBalanceUseCase(balanceusecase.New(store, log)),
	).Configure()

	log.WithFields(logrus.Fields{
		"address": *cfg.RunAddress,
	}).Info("start server")

	err = http.ListenAndServe(*cfg.RunAddress, rtr)
	if err != nil {
		return err
	}

	return nil
}

func shutdownHandler(cancel context.CancelFunc, log *logger.Logger) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		log.WithFields(logrus.Fields{
			"signal": sig,
		}).Info("kill process")

		cancel()

		os.Exit(0)
	}()
}
