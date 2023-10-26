package accrualerrors

type RateLimitError struct {
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return "count reqs is many: over rate limit app"
}
