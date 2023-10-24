package entities

import (
	"github.com/MaximPolyaev/gofermart/internal/enums/accrualstatus"
)

type AccrualOrder struct {
	Order   string                      `json:"order"`
	Status  accrualstatus.AccrualStatus `json:"status"`
	Accrual float64                     `json:"accrual"`
}
