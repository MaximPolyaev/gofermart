package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Http struct {
	runAddr string
}

func New(runAddr string) *Http {
	return &Http{
		runAddr: runAddr,
	}
}

func (a *Http) ListenAndServe(router *chi.Mux) error {
	return http.ListenAndServe(a.runAddr, router)
}
