package db_utils

import (
	"context"
	"log"
	
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

func InsertReservationData(ctx context.Context, conn *pgx.Conn, r models.Reservation) (*models.Reservation, error) {
	args := pgx.NamedArgs{
		"reservation_kind":    r.Kind,
		"tunnel_id":           r.TunnelId,
		"coach_id":            r.CoachId,
		"customer_first_name": r.CustomerFirstName,
		"customer_last_name":  r.CustomerLastName,
		"customer_phone":      r.CustomerPhone,
		"customer_email":      r.CustomerEmail,
		"start_time":          r.StartTime,
		"duration_minutes":    r.Duration,
		"end_time":            r.EndTime,
		"status":              r.Status,
		"notes":               r.Notes,
	}
	
	const query = `
		INSERT INTO reservations (
			reservation_kind,
			tunnel_id,
			coach_id,
			customer_first_name,
			customer_last_name,
			customer_phone,
			customer_email,
			start_time,
			duration_minutes,
			end_time,
			status,
			notes
		)

		VALUES (
			@reservation_kind,
			@tunnel_id,
			@coach_id,
			@customer_first_name,
			@customer_last_name,
			@customer_phone,
			@customer_email,
			@start_time,
			@duration_minutes,
			@end_time,
			@status,
			@notes
		)

		RETURNING *;
	`
	
	var out models.Reservation
	if err := pgxscan.Get(ctx, conn, &out, query, args); err != nil {
		return nil, err
	}
	
	return &out, nil
}

func UpdateReservationData(ctx context.Context, conn *pgx.Conn, id string, reservation models.ReservationUpdate) (*models.Reservation, error) {
	var updatedReservation models.Reservation
	
	// send update to the db
	args := pgx.NamedArgs{
		"id":                  id,
		"reservation_kind":    reservation.Kind,
		"tunnel_id":           reservation.TunnelId,
		"coach_id":            reservation.CoachId,
		"customer_first_name": reservation.CustomerFirstName,
		"customer_last_name":  reservation.CustomerLastName,
		"customer_phone":      reservation.CustomerPhone,
		"customer_email":      reservation.CustomerEmail,
		"start_time":          reservation.StartTime,
		"duration_minutes":    reservation.Duration,
		"end_time":            reservation.EndTime,
		"status":              reservation.Status,
		"notes":               reservation.Notes,
	}
	
	query := `
			UPDATE reservations
			SET
				reservation_kind = COALESCE(@reservation_kind, reservation_kind),
				tunnel_id = COALESCE(@tunnel_id, tunnel_id),
				coach_id = COALESCE(@coach_id, coach_id),
				customer_first_name = COALESCE(@customer_first_name, customer_first_name),
				customer_last_name = COALESCE(@customer_last_name, customer_last_name),
				customer_phone = COALESCE(@customer_phone, customer_phone),
				customer_email = COALESCE(@customer_email, customer_email),
				start_time = COALESCE(@start_time, start_time),
				duration_minutes = COALESCE(@duration_minutes, duration_minutes),
				end_time = COALESCE(@end_time, end_time),
				status = COALESCE(@status, status),
				notes = COALESCE(@notes, notes),
				updated_at = now()
			WHERE id = @id
			RETURNING *
		`
	
	err := pgxscan.Get(ctx, conn, &updatedReservation, query, args)
	
	if pgxscan.NotFound(err) {
		log.Println("[API] Could not find reservation with id:", id)
		return nil, nil
	}
	
	if err != nil {
		log.Println("[API] Error updating reservation:", err)
		return nil, err
	}
	
	return &updatedReservation, nil
}

func LoadReservationData(ctx context.Context, conn *pgx.Conn) ([]models.Reservation, error) {
	reservations := make([]models.Reservation, 0)
	
	err := pgxscan.Select(
		ctx,
		conn,
		&reservations,
		`SELECT * FROM reservations`,
	)
	
	if err != nil {
		log.Println("[API] Error querying database:", err)
		return nil, err
	}
	
	return reservations, nil
}

func LoadReservationDataWithParams(ctx context.Context, conn *pgx.Conn, fromTime, toTime, tunnelId string) ([]models.Reservation, error) {
	// create empty slice so it doesn't respond with nil
	reservations := make([]models.Reservation, 0)
	
	args := pgx.NamedArgs{
		"from_time": fromTime,
		"to_time":   toTime,
		"tunnel_id": tunnelId,
	}
	
	var query string
	if tunnelId == "" {
		query = `
			SELECT * FROM reservations
			WHERE start_time >= @from_time::timestamptz AND start_time < @to_time::timestamptz
			ORDER BY start_time ASC
		`
	} else {
		query = `
			SELECT * FROM reservations
			WHERE tunnel_id = @tunnel_id::int
				AND start_time >= @from_time::timestamptz
				AND start_time < @to_time::timestamptz
			ORDER BY start_time ASC
		`
	}
	
	err := pgxscan.Select(ctx, conn, &reservations, query, args)
	if err != nil {
		log.Println("[API] Error querying database:", err)
		return nil, err
	}
	
	if len(reservations) == 0 {
		log.Println("[API] No reservations found that matched the search params -",
			"fromTime:",
			fromTime,
			"toTime:",
			toTime,
			"tunnelId:",
			tunnelId,
		)
	}
	
	return reservations, nil
}

func LoadReservationById(ctx context.Context, conn *pgx.Conn, id string) (*models.Reservation, error) {
	var reservation models.Reservation
	
	query := `SELECT * FROM reservations WHERE id=$1`
	
	err := pgxscan.Get(
		ctx,
		conn,
		&reservation,
		query,
		id,
	)
	
	if pgxscan.NotFound(err) {
		log.Println("[API] No rows found while querying reservations")
		return nil, nil
	} else if err != nil {
		log.Println("[API] Error querying database:", err)
		return nil, err
	}
	
	return &reservation, nil
}

func DeleteReservationData(ctx context.Context, conn *pgx.Conn, id string) (int64, error) {
	cmdTag, err := conn.Exec(
		ctx,
		"DELETE FROM reservations WHERE id=$1",
		id,
	)
	
	if err != nil {
		log.Println("[API] Error deleting reservation:", err)
		return 0, err
	}
	
	return cmdTag.RowsAffected(), err
}
