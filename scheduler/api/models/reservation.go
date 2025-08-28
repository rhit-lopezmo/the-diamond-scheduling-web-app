package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ReservationKind string

const (
	ReservationKindTunnel ReservationKind = "tunnel"
	ReservationKindLesson ReservationKind = "lesson"
)

type ReservationStatus string

const (
	ReservationStatusHeld      ReservationStatus = "held"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusCancelled ReservationStatus = "cancelled"
	ReservationStatusCompleted ReservationStatus = "completed"
	ReservationStatusNoShow    ReservationStatus = "no_show"
)

type Reservation struct {
	Id                pgtype.UUID        `db:"id" json:"id"`
	Kind              ReservationKind    `db:"reservation_kind" json:"reservation_kind"`
	TunnelId          *int32             `db:"tunnel_id" json:"tunnel_id"`
	CoachId           *pgtype.UUID       `db:"coach_id" json:"coach_id"`
	CustomerFirstName string             `db:"customer_first_name" json:"customer_first_name"`
	CustomerLastName  string             `db:"customer_last_name" json:"customer_last_name"`
	CustomerPhone     string             `db:"customer_phone" json:"customer_phone"`
	CustomerEmail     *string            `db:"customer_email" json:"customer_email"`
	StartTime         pgtype.Timestamptz `db:"start_time" json:"start_time"`
	Duration          int32              `db:"duration_minutes" json:"duration_minutes"`
	EndTime           pgtype.Timestamptz `db:"end_time" json:"end_time"`
	Status            ReservationStatus  `db:"status" json:"status"`
	Notes             *string            `db:"notes" json:"notes"`
	CreatedAt         pgtype.Timestamptz `db:"created_at" json:"created_at"`
}
