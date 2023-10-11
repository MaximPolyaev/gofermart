package router

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	*chi.Mux
	auth AuthUseCase
}

func New(auth AuthUseCase) *Router {
	router := chi.NewRouter()

	return &Router{Mux: router, auth: auth}
}

func (r *Router) Configure() {
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(gzip.BestCompression))

	r.Post("/api/user/login", r.login())
	r.Post("/api/user/register", r.register())
	r.Post("/api/user/orders", r.postOrders())
	r.Get("/api/user/orders", r.getOrders())
	r.Get("/api/user/balance", r.balance())
	r.Post("/api/user/balance/withdraw", r.withdraw())
	r.Get("/api/user/withdrawals", r.withdrawInfo())
}

func (r *Router) postOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
}

func (r *Router) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
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
