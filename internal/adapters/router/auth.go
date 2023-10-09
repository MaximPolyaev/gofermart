package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/entities"
)

type AuthUseCase interface {
	ValidatePayload(payload entities.AuthPayload) error
	SignIn(ctx context.Context, payload entities.AuthPayload) (string, error)
	SignUp(ctx context.Context, payload entities.AuthPayload) (string, error)
}

func (r *Router) login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ct := req.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "unexpected content type "+ct, http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := req.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var payload entities.AuthPayload

		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := r.auth.ValidatePayload(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := r.auth.SignIn(req.Context(), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r.writeAuthToken(w.Header(), token)

		w.WriteHeader(http.StatusOK)
	}
}

func (r *Router) registration() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ct := req.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "unexpected content type "+ct, http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := req.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var payload entities.AuthPayload

		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := r.auth.ValidatePayload(payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := r.auth.SignUp(req.Context(), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		r.writeAuthToken(w.Header(), token)

		w.WriteHeader(http.StatusOK)
	}
}
