package notifications

import (
	"github.com/izaakdale/service-event/pkg/proto/event"
)

type OrderCreatedPayload struct {
	OrderID      string              `json:"order_id,omitempty"`
	OrderRequest *event.OrderRequest `json:"order,omitempty"`
}
