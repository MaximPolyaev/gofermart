package router

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/MaximPolyaev/gofermart/internal/errors/balanceerrors"
	"io"
	"net/http"
)

type balanceUseCase interface {
	GetBalance(ctx context.Context, userID int) (*entities.UserBalance, error)
	WriteOff(ctx context.Context, off entities.WriteOff, userID int) error
	GetWroteOffs(ctx context.Context, userID int) ([]entities.WroteOff, error)
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
	return func(w http.ResponseWriter, req *http.Request) {
		ct := req.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "unexpected content type "+ct, http.StatusBadRequest)
			return
		}

		data, err := io.ReadAll(req.Body)
		defer func() {
			_ = req.Body.Close()
		}()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var writeOff entities.WriteOff
		err = json.Unmarshal(data, &writeOff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if writeOff.Sum <= 0 {
			http.Error(w, "incorrect write off sum", http.StatusBadRequest)
			return
		}

		userID, err := r.getUserIDFromReq(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		err = r.orders.ValidateNumber(writeOff.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		rctx := req.Context()

		userIDByOrder, err := r.orders.GetUserID(rctx, writeOff.Order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userIDByOrder == 0 {
			err = r.orders.CreateOrder(rctx, writeOff.Order, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if userIDByOrder != userID {
			http.Error(w, "order by is assignment another user", http.StatusUnprocessableEntity)
			return
		}

		err = r.balance.WriteOff(rctx, writeOff, userID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, balanceerrors.ErrInsufficientFunds) {
				statusCode = http.StatusPaymentRequired
			}

			http.Error(w, err.Error(), statusCode)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (r *Router) withdrawInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, err := r.getUserIDFromReq(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		wroteOffs, err := r.balance.GetWroteOffs(req.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(wroteOffs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		data, err := json.Marshal(wroteOffs)
		if err == nil {
			_, err = w.Write(data)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
