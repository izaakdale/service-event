// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"time"
)

type Event struct {
	EventID          int64     `json:"event_id"`
	EventName        string    `json:"event_name"`
	TicketsRemaining int32     `json:"tickets_remaining"`
	EventTimestamp   time.Time `json:"event_timestamp"`
}