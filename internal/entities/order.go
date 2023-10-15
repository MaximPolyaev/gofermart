package entities

import (
	"time"

	"github.com/MaximPolyaev/gofermart/internal/enums"
)

type RFC3339Time time.Time

type Order struct {
	Number     string            `json:"number"`
	Status     enums.OrderStatus `json:"status"`
	Accrual    float64           `json:"accrual,omitempty"`
	UploadedAt RFC3339Time       `json:"uploaded_at"`
}

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte((time.Time(t)).Format("\"" + time.RFC3339 + "\"")), nil
}
