package models

import "github.com/jackc/pgx/v5/pgtype"

type Tunnel struct {
	Id        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	IsActive  bool               `json:"is_active"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
