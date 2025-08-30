package models

import "github.com/jackc/pgx/v5/pgtype"

const (
	SpecialtyHitting  string = "hitting"
	SpecialtyPitching string = "pitching"
	SpecialtyFielding string = "fielding"
	SpecialtyCatching string = "catching"
)

type Coach struct {
	Id          pgtype.UUID        `db:"id" json:"id"`
	FirstName   string             `db:"first_name" json:"first_name"`
	LastName    string             `db:"last_name" json:"last_name"`
	Email       *string            `db:"email" json:"email"`
	Phone       string             `db:"phone" json:"phone"`
	IsActive    bool               `db:"is_active" json:"is_active"`
	Specialties []string           `db:"specialties" json:"specialties"`
	CreatedAt   pgtype.Timestamptz `db:"created_at" json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `db:"updated_at" json:"updated_at"`
}

type CoachUpdates struct {
	FirstName   *string   `db:"first_name" json:"first_name"`
	LastName    *string   `db:"last_name" json:"last_name"`
	Email       *string   `db:"email" json:"email"`
	Phone       *string   `db:"phone" json:"phone"`
	IsActive    *bool     `db:"is_active" json:"is_active"`
	Specialties *[]string `db:"specialties" json:"specialties"`
}
