package entities

type WriteOff struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WroteOff struct {
	WriteOff
	ProcessedAt RFC3339Time `json:"processed_at"`
}
