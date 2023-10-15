package router

import (
	"context"
	"io"
	"net/http"
	"strconv"
)

type ordersUseCase interface {
	ValidateLuhn(number int64) bool
	GetUserID(ctx context.Context, number int64) (int, error)
	CreateOrder(ctx context.Context, number int64, userID int) error
}

func (r *Router) postOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		ct := req.Header.Get("Content-Type")
		if ct != "text/plain" {
			http.Error(w, "incorrect format request", http.StatusBadRequest)
			return
		}

		number, err := r.getOrderNumberFromReq(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !r.orders.ValidateLuhn(number) {
			http.Error(w, "incorrect number format", http.StatusUnprocessableEntity)
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

func (r *Router) getOrderNumberFromReq(req *http.Request) (int64, error) {
	data, err := io.ReadAll(req.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)

	if err != nil {
		return 0, err
	}

	number, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
