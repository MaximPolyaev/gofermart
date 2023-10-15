package router

import (
	"context"
	"io"
	"net/http"
)

type ordersUseCase interface {
	ValidateNumber(number string) error
	GetUserID(ctx context.Context, number string) (int, error)
	CreateOrder(ctx context.Context, number string, userID int) error
}

func (r *Router) postOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		ct := req.Header.Get("Content-Type")
		if ct != "text/plain" {
			http.Error(w, "incorrect format request", http.StatusBadRequest)
			return
		}

		data, err := io.ReadAll(req.Body)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		number := string(data)

		if err := r.orders.ValidateNumber(number); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		rctx := req.Context()

		userID, err := r.getUserIDFromReq(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userIDByExistOrder, err := r.orders.GetUserID(rctx, number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userIDByExistOrder != 0 {
			if userIDByExistOrder == userID {
				w.WriteHeader(http.StatusOK)
				return
			}

			http.Error(w, "order will be loaded on another user", http.StatusConflict)
			return
		}

		err = r.orders.CreateOrder(rctx, number, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func (r *Router) getOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}
}
