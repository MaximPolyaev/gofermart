package router

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MaximPolyaev/gofermart/internal/utils/jwt"
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

func (r *Router) getClaimsFromReq(req *http.Request) (*jwt.Claims, error) {
	token := r.getToken(req)

	if token == "" {
		return nil, errors.New("token could not be determined")
	}

	claims, err := jwt.ParseClaims(token)

	if err != nil {
		return nil, err
	}

	return claims, nil
}
