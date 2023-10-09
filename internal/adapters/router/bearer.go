package router

import "net/http"

const (
	authKey   = "Authorization"
	bearerKey = "Bearer"
)

func (r *Router) writeAuthToken(header http.Header, token string) {
	header.Set(authKey, bearerKey+" "+token)
}
