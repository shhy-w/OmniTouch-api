package runtime_settings

import (
	"encoding/json"
	"github.com/uozi-tech/cosy/logger"
)

type Order struct {
	PendingPaymentTimeout uint `json:"pending_payment_timeout"` // day(s)
	AfterSaleTimeout      uint `json:"after_sale_timeout"`      // day(s)
}

// GetPendingPaymentTimeout expired after order created (days)
func (o *Order) GetPendingPaymentTimeout() uint {
	if o.PendingPaymentTimeout == 0 {
		return 15
	}
	return o.PendingPaymentTimeout
}

// GetAfterSaleTimeout expired after order paid (days)
func (o *Order) GetAfterSaleTimeout() uint {
	if o.AfterSaleTimeout == 0 {
		return 60
	}
	return o.AfterSaleTimeout
}

func GetOrderSettings() (o *Order) {
	o = &Order{}
	err := Get("order", o)
	if err != nil {
		logger.Error(err)
	}
	return
}

func UpdateOrderSettings(target *Order) {
	bytes, err := json.Marshal(target)
	if err != nil {
		logger.Error(err)
		return
	}
	Set("order", string(bytes))
}
