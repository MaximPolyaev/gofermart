package entities

import "github.com/MaximPolyaev/gofermart/internal/enums/orderstatus"

type AccrualOrder struct {
	Order   string                  `json:"order"`
	Status  orderstatus.OrderStatus `json:"status"`
	Accrual float64                 `json:"accrual"`
}
