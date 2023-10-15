package entities

type WriteOff struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
