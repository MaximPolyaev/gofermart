package accrualerrors

import "errors"

var ErrRateLimit = errors.New("count reqs is many: over rate limit app")
