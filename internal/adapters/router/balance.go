package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type balanceUseCase interface {
	GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error)
}

func WithBalanceUseCase(balance balanceUseCase) func(r *Router) {
	return func(r *Router) {
		r.balance = balance
	}
}

func (r *Router) getBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, err := r.getUserIDFromReq(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		fmt.Println(userID)

		balance, err := r.balance.GetBalance(req.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(*balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
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
