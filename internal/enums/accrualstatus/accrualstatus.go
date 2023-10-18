package accrualstatus

type AccrualStatus string

const (
	REGISTERED = AccrualStatus("REGISTERED")
	INVALID    = AccrualStatus("INVALID")
	PROCESSING = AccrualStatus("PROCESSING")
	PROCESSED  = AccrualStatus("PROCESSED")
)
