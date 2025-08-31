package db_utils

import (
	"context"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

func LoadCoachesData(ctx context.Context, conn *pgx.Conn) ([]models.Coach, error) {
	coaches := make([]models.Coach, 0)

	query := `SELECT * FROM coaches`

	err := pgxscan.Select(ctx, conn, &coaches, query)
	if err != nil {
		log.Println("[API] Error querying database:", err)
		return nil, err
	}

	return coaches, nil
}

func UpdateCoachData(ctx context.Context, conn *pgx.Conn, id string, updates models.CoachUpdates) (*models.Coach, error) {
	var updatedCoach models.Coach

	// send update to the db
	args := pgx.NamedArgs{
		"id":          id,
		"first_name":  updates.FirstName,
		"last_name":   updates.LastName,
		"phone":       updates.Phone,
		"email":       updates.Email,
		"is_active":   updates.IsActive,
		"specialties": updates.Specialties,
	}

	query := `
			UPDATE coaches
			SET
				first_name = COALESCE(@first_name, first_name),
				last_name = COALESCE(@last_name, last_name),
				phone = COALESCE(@phone, phone),
				email = COALESCE(@email, email),
				is_active = COALESCE(@is_active, is_active),
				specialties = COALESCE(@specialties, specialties),
				updated_at = now()
			WHERE id = @id
			RETURNING *
		`

	err := pgxscan.Get(ctx, conn, &updatedCoach, query, args)

	if pgxscan.NotFound(err) {
		log.Println("[API] Could not find reservation with id:", id)
		return nil, nil
	}

	if err != nil {
		log.Println("[API] Error updating reservation:", err)
		return nil, err
	}

	return &updatedCoach, nil
}

func InsertCoachData(ctx context.Context, conn *pgx.Conn, c models.Coach) (*models.Coach, error) {
	args := pgx.NamedArgs{
		"first_name":  c.FirstName,
		"last_name":   c.LastName,
		"phone":       c.Phone,
		"email":       c.Email,
		"specialties": c.Specialties,
	}

	const query = `
		INSERT INTO coaches (
			first_name,
			last_name,
			phone,
			email,
			specialties
		)

		VALUES (
			@first_name,
			@last_name,
			@phone,
			@email,
			@specialties::coach_specialty[]
		)

		RETURNING *;
	`

	var out models.Coach
	if err := pgxscan.Get(ctx, conn, &out, query, args); err != nil {
		return nil, err
	}

	return &out, nil
}
