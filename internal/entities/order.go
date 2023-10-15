package entities

import (
	"time"

	"github.com/MaximPolyaev/gofermart/internal/enums"
)

type Order struct {
	Number     int
	Status     enums.OrderStatus
	Accrual    int
	UploadedAt time.Time
}
