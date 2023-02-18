package app

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/izaakdale/lib/publisher"
	db "github.com/izaakdale/service-event/internal/datastore/sqlc"
	"github.com/izaakdale/service-event/pkg/notifications"
	"github.com/izaakdale/service-event/pkg/proto/event"
	_ "github.com/lib/pq"
)

func (g *GServer) GetEvent(ctx context.Context, ev *event.EventRequest) (*event.EventResponse, error) {
	dbe, err := querier.GetEvent(ctx, ev.EventId)
	if err != nil {
		return nil, err
	}
	return &event.EventResponse{
		EventId:          dbe.EventID,
		EventName:        dbe.EventName,
		TicketsRemaining: dbe.TicketsRemaining,
		Timestamp:        dbe.EventTimestamp.Format(time.RFC3339),
	}, nil
}

func (g *GServer) GetEvents(ctx context.Context, le *event.ListEventRequest) (*event.ListEventResponse, error) {
	dbes, err := querier.GetEvents(ctx, le.EventIds)
	if err != nil {
		return nil, err
	}
	ret := &event.ListEventResponse{}
	for _, e := range dbes {
		ret.Events = append(ret.Events, &event.EventResponse{
			EventId:          e.EventID,
			EventName:        e.EventName,
			TicketsRemaining: e.TicketsRemaining,
			Timestamp:        e.EventTimestamp.Format(time.RFC3339),
		})
	}
	return ret, nil
}

func (g *GServer) MakeOrder(ctx context.Context, e *event.OrderRequest) (*event.OrderResponse, error) {
	// create an order UUID
	id := uuid.New().String()
	log.Printf("placing order %s\n", id)

	dbe, err := querier.UpdateEvent(ctx, db.UpdateEventParams{
		EventID:  e.EventId,
		NTickets: int32(len(e.Attendees)),
	})
	if err != nil {
		return nil, err
	}

	// publish to SNS
	payload := notifications.OrderCreatedPayload{
		OrderID:      id,
		OrderRequest: e,
	}

	// not using message id for now.
	oBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	_, err = publisher.Publish(string(oBytes))
	if err != nil {
		return nil, err
	}

	return &event.OrderResponse{
		EventId: dbe.EventID,
		OrderId: id,
	}, nil
}
