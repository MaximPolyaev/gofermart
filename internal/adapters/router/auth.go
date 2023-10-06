package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type AuthUseCase interface {
	signIn(ctx context.Context, payload entities.AuthPayload) (string, error)
	signUp(ctx context.Context, payload entities.AuthPayload) (string, error)
}

func (r *Router) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "unexpected content type "+ct, http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := r.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var payload entities.AuthPayload

		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (r *Router) registration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "unexpected content type "+ct, http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := r.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var payload entities.AuthPayload

		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
