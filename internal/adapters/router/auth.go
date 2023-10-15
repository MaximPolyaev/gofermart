package router

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MaximPolyaev/gofermart/internal/entities"
	"github.com/go-chi/chi/v5"
)

type authUseCase interface {
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

func (r *Router) register() http.HandlerFunc {
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

func (r *Router) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if r.isNeedCheckAuth(req) && !r.isAuthenticated(req) {
			http.Error(w, "no authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *Router) isNeedCheckAuth(req *http.Request) bool {
	routePath := req.URL.Path

	rctx := chi.RouteContext(req.Context())

	if rctx.RoutePath != "" {
		routePath = rctx.RoutePath
	}

	return routePath != "/api/user/login" && routePath != "/api/user/register"
}

func (r *Router) isAuthenticated(req *http.Request) bool {
	claims, err := r.getClaimsFromReq(req)

	return err == nil && claims != nil
}
