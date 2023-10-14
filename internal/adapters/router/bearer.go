package router

import (
	"net/http"
	"strings"
)

const (
	authKey   = "Authorization"
	bearerKey = "Bearer"
)

func (r *Router) writeAuthToken(header http.Header, token string) {
	header.Set(authKey, bearerKey+" "+token)
}

func (r *Router) getToken(req *http.Request) string {
	authHeader := req.Header.Get(authKey)

	if authHeader == "" {
		return ""
	}

	return strings.Replace(authHeader, bearerKey+" ", "", 1)
}
