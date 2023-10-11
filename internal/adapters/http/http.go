package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	runAddr string
}

func New(runAddr string) *Server {
	return &Server{
		runAddr: runAddr,
	}
}

func (a *Server) ListenAndServe(router *chi.Mux) error {
	fmt.Println("start server on", a.runAddr)
	return http.ListenAndServe(a.runAddr, router)
}
