package router

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
	auth   authUseCase
	orders ordersUseCase
	user   userUseCase
}

func New(useCases ...func(r *Router)) *Router {
	mux := chi.NewRouter()

	router := &Router{Mux: mux}

	for _, uc := range useCases {
		uc(router)
	}

	return router
}

func WithAuthUseCase(auth authUseCase) func(r *Router) {
	return func(r *Router) {
		r.auth = auth
	}
}

func WithOrdersUseCase(orders ordersUseCase) func(r *Router) {
	return func(r *Router) {
		r.orders = orders
	}
}

func WithUserUseCase(user userUseCase) func(r *Router) {
	return func(r *Router) {
		r.user = user
	}
}

func (r *Router) Configure() {
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
	r.Get("/api/user/balance", r.balance())
	r.Post("/api/user/balance/withdraw", r.withdraw())
	r.Get("/api/user/withdrawals", r.withdrawInfo())
}

func (r *Router) balance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
}

func (r *Router) withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
}

func (r *Router) withdrawInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
}
