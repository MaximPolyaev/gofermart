package entities

import (
	"github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"
)

type Order struct {
	Number     string                  `json:"number"`
	Status     orderstatus.OrderStatus `json:"status"`
	Accrual    float64                 `json:"accrual,omitempty"`
	UploadedAt RFC3339Time             `json:"uploaded_at"`
}
