package router

import (
	"compress/gzip"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
	log     logger
	auth    authUseCase
	orders  ordersUseCase
	user    userUseCase
	balance balanceUseCase
}

type logger interface {
	Info(args ...interface{})
}

func New(log logger, useCases ...func(r *Router)) *Router {
	mux := chi.NewRouter()

	router := &Router{
		Mux: mux,
		log: log,
	}

	for _, uc := range useCases {
		uc(router)
	}

	return router
}

func (r *Router) Configure() *Router {
	r.Use(
		middleware.Logger,
		middleware.StripSlashes,
		r.authMiddleware,
		middleware.Compress(gzip.BestCompression),
	)

	r.Post("/api/user/login", r.login())
	r.Post("/api/user/register", r.register())
	r.Post("/api/user/orders", r.postOrders())
	r.Get("/api/user/orders", r.getOrders())
	r.Get("/api/user/balance", r.getBalance())
	r.Post("/api/user/balance/withdraw", r.withdraw())
	r.Get("/api/user/withdrawals", r.withdrawInfo())

	return r
}
