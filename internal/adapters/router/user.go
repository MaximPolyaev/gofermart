package router

import (
	"context"
	"net/http"
)

type userUseCase interface {
	GetUserIDByLogin(ctx context.Context, login string) (int, error)
}

func (r *Router) getUserIDFromReq(req *http.Request) (int, error) {
	claims, err := r.getClaimsFromReq(req)
	if err != nil {
		return 0, err
	}

	id, err := r.user.GetUserIDByLogin(req.Context(), claims.UserLogin)
	if err != nil {
		return 0, err
	}

	return id, nil
}
