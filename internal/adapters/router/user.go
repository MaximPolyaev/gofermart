package router

import (
	"context"
	"net/http"
)

type userUseCase interface {
	GetUserIdByLogin(ctx context.Context, login string) (int, error)
}

func (r *Router) getUserIdFromReq(req *http.Request) (int, error) {
	claims, err := r.getClaimsFromReq(req)
	if err != nil {
		return 0, err
	}

	id, err := r.user.GetUserIdByLogin(req.Context(), claims.UserLogin)
	if err != nil {
		return 0, err
	}

	return id, nil
}
