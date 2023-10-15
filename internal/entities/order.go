package entities

import (
	"github.com/MaximPolyaev/gofermart/internal/enums"
)

type Order struct {
	Number     string            `json:"number"`
	Status     enums.OrderStatus `json:"status"`
	Accrual    float64           `json:"accrual,omitempty"`
	UploadedAt RFC3339Time       `json:"uploaded_at"`
}
