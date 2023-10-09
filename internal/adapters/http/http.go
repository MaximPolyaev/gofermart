package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	fmt.Println("start server on", a.runAddr)
	return http.ListenAndServe(a.runAddr, router)
}
