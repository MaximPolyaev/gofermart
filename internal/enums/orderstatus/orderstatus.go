package orderstatus

type OrderStatus string

const (
	NEW        = OrderStatus("NEW")
	PROCESSING = OrderStatus("PROCESSING")
	INVALID    = OrderStatus("INVALID")
	PROCESSED  = OrderStatus("PROCESSED")
)
