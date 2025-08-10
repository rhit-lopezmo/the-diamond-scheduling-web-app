package models

import "time"

type Reservation struct {
	Id         string    `json:"id"`
	TunnelId   int32     `json:"tunnel_id"`
	CustomerId string    `json:"customer_id"`
	Title      string    `json:"title"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
	Notes      string    `json:"notes"`
}
