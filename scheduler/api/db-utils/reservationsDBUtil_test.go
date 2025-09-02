package db_utils

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

func Test_LoadReservationsData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	u2 := uuid.New()
	pgU1 := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}
	pgU2 := pgtype.UUID{Bytes: [16]byte(u2), Valid: true}

	rows := pgxmock.NewRows([]string{"id", "customer_first_name", "customer_last_name"}).
		AddRow(pgU1, "John", "Doe").
		AddRow(pgU2, "Jane", "Doe")

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM reservations`)).WillReturnRows(rows)

	// exercise
	result, err := LoadReservationData(context.Background(), mockConn)

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(result) != 2 {
		t.Fatal("expected 2 rows, got", len(result))
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadReservationsData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM reservations`)).WillReturnError(errors.New("test error"))

	// exercise
	result, err := LoadReservationData(context.Background(), mockConn)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if result != nil {
		t.Fatal("expected no result, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadReservationsDataById(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	pgU1 := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}

	rows := pgxmock.NewRows([]string{"id", "customer_first_name", "customer_last_name"}).
		AddRow(pgU1, "John", "Doe")

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM reservations WHERE id=$1`)).
		WithArgs(pgU1.String()).
		WillReturnRows(rows)

	// exercise
	result, err := LoadReservationById(context.Background(), mockConn, pgU1.String())

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result == nil {
		t.Fatal("expected result, got none")
	}

	if result.Id != pgU1 {
		t.Fatal("expected", pgU1, "got", result.Id)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadReservationsDataById_NotFound(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	pgU1 := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM reservations WHERE id=$1`)).
		WithArgs(pgU1.String()).
		WillReturnError(pgx.ErrNoRows)

	// exercise
	result, err := LoadReservationById(context.Background(), mockConn, pgU1.String())

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result != nil {
		t.Fatal("expected no result, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadReservationsDataById_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	pgU1 := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM reservations WHERE id=$1`)).
		WithArgs(pgU1.String()).
		WillReturnError(errors.New("test error"))

	// exercise
	result, err := LoadReservationById(context.Background(), mockConn, pgU1.String())

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if result != nil {
		t.Fatal("expected no result, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadReservationsDataWithParams(t *testing.T) {

}

func Test_LoadReservationsDataWithParams_Error(t *testing.T) {

}

func Test_InsertReservationData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	query := `
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
	u := uuid.New()
	pgU := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	var tunnelId int32 = 1

	testReservation := models.Reservation{
		Id:                pgU,
		Kind:              models.ReservationKindTunnel,
		TunnelId:          &tunnelId,
		CoachId:           nil,
		CustomerFirstName: "John",
		CustomerLastName:  "Doe",
		CustomerPhone:     "1112223333",
		CustomerEmail:     nil,
		StartTime:         pgtype.Timestamptz{},
		Duration:          0,
		EndTime:           pgtype.Timestamptz{},
		Status:            models.ReservationStatusConfirmed,
		Notes:             nil,
	}

	rows := pgxmock.NewRows([]string{"id", "customer_first_name", "customer_last_name"}).
		AddRow(
			testReservation.Id,
			testReservation.CustomerFirstName,
			testReservation.CustomerLastName,
		)

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pgx.NamedArgs{
		"reservation_kind":    testReservation.Kind,
		"tunnel_id":           testReservation.TunnelId,
		"coach_id":            testReservation.CoachId,
		"customer_first_name": testReservation.CustomerFirstName,
		"customer_last_name":  testReservation.CustomerLastName,
		"customer_phone":      testReservation.CustomerPhone,
		"customer_email":      testReservation.CustomerEmail,
		"start_time":          testReservation.StartTime,
		"duration_minutes":    testReservation.Duration,
		"end_time":            testReservation.EndTime,
		"status":              testReservation.Status,
		"notes":               testReservation.Notes,
	}).WillReturnRows(rows)

	// exercise
	result, err := InsertReservationData(context.Background(), mockConn, testReservation)

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result == nil {
		t.Fatal("expected result, got none")
	}

	if result.Id != testReservation.Id {
		t.Fatal("expected id ", testReservation.Id.String(), "got", result.Id.String())
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_InsertReservationData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	query := `
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
	u := uuid.New()
	pgU := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	var tunnelId int32 = 1

	testReservation := models.Reservation{
		Id:                pgU,
		Kind:              models.ReservationKindTunnel,
		TunnelId:          &tunnelId,
		CoachId:           nil,
		CustomerFirstName: "John",
		CustomerLastName:  "Doe",
		CustomerPhone:     "1112223333",
		CustomerEmail:     nil,
		StartTime:         pgtype.Timestamptz{},
		Duration:          0,
		EndTime:           pgtype.Timestamptz{},
		Status:            models.ReservationStatusConfirmed,
		Notes:             nil,
	}

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pgx.NamedArgs{
		"reservation_kind":    testReservation.Kind,
		"tunnel_id":           testReservation.TunnelId,
		"coach_id":            testReservation.CoachId,
		"customer_first_name": testReservation.CustomerFirstName,
		"customer_last_name":  testReservation.CustomerLastName,
		"customer_phone":      testReservation.CustomerPhone,
		"customer_email":      testReservation.CustomerEmail,
		"start_time":          testReservation.StartTime,
		"duration_minutes":    testReservation.Duration,
		"end_time":            testReservation.EndTime,
		"status":              testReservation.Status,
		"notes":               testReservation.Notes,
	}).WillReturnError(errors.New("test error"))

	// exercise
	result, err := InsertReservationData(context.Background(), mockConn, testReservation)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if result != nil {
		t.Fatal("expected no result, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_UpdateReservationData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

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

	firstName := "John"
	lastName := "Doe"
	phone := "1112223333"

	u := uuid.New()
	pgUUID := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	testReservationUpdates := models.ReservationUpdates{
		CustomerFirstName: &firstName,
		CustomerLastName:  &lastName,
		CustomerPhone:     &phone,
	}

	rows := pgxmock.NewRows([]string{"id", "customer_first_name", "customer_last_name", "customer_phone"}).
		AddRow(
			pgUUID,
			*testReservationUpdates.CustomerFirstName,
			*testReservationUpdates.CustomerLastName,
			*testReservationUpdates.CustomerPhone,
		)

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pgx.NamedArgs{
		"reservation_kind":    testReservationUpdates.Kind,
		"tunnel_id":           testReservationUpdates.TunnelId,
		"coach_id":            testReservationUpdates.CoachId,
		"customer_first_name": testReservationUpdates.CustomerFirstName,
		"customer_last_name":  testReservationUpdates.CustomerLastName,
		"customer_phone":      testReservationUpdates.CustomerPhone,
		"customer_email":      testReservationUpdates.CustomerEmail,
		"start_time":          testReservationUpdates.StartTime,
		"duration_minutes":    testReservationUpdates.Duration,
		"end_time":            testReservationUpdates.EndTime,
		"status":              testReservationUpdates.Status,
		"notes":               testReservationUpdates.Notes,
		"id":                  pgUUID.String(),
	}).WillReturnRows(rows)

	// exercise
	result, err := UpdateReservationData(context.Background(), mockConn, pgUUID.String(), testReservationUpdates)

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result == nil {
		t.Fatal("expected result, got none")
	}

	if result.Id != pgUUID {
		t.Fatal("expected id ", pgUUID.String(), "got", result.Id.String())
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_UpdateReservationData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

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

	firstName := "John"
	lastName := "Doe"
	phone := "1112223333"

	u := uuid.New()
	pgUUID := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	testReservationUpdates := models.ReservationUpdates{
		CustomerFirstName: &firstName,
		CustomerLastName:  &lastName,
		CustomerPhone:     &phone,
	}

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(pgx.NamedArgs{
		"reservation_kind":    testReservationUpdates.Kind,
		"tunnel_id":           testReservationUpdates.TunnelId,
		"coach_id":            testReservationUpdates.CoachId,
		"customer_first_name": testReservationUpdates.CustomerFirstName,
		"customer_last_name":  testReservationUpdates.CustomerLastName,
		"customer_phone":      testReservationUpdates.CustomerPhone,
		"customer_email":      testReservationUpdates.CustomerEmail,
		"start_time":          testReservationUpdates.StartTime,
		"duration_minutes":    testReservationUpdates.Duration,
		"end_time":            testReservationUpdates.EndTime,
		"status":              testReservationUpdates.Status,
		"notes":               testReservationUpdates.Notes,
		"id":                  pgUUID.String(),
	}).WillReturnError(errors.New("test error"))

	// exercise
	result, err := UpdateReservationData(context.Background(), mockConn, pgUUID.String(), testReservationUpdates)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if result != nil {
		t.Fatal("expected no result, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteReservationData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	pgUUID := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}

	mockConn.ExpectExec(regexp.QuoteMeta(`DELETE FROM reservations WHERE id=$1`)).
		WithArgs(pgUUID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// exercise
	result, err := DeleteReservationData(context.Background(), mockConn, pgUUID.String())

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result != 1 {
		t.Fatal("expected 1 deleted row, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteReservationData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	pgUUID := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}

	mockConn.ExpectExec(regexp.QuoteMeta(`DELETE FROM reservations WHERE id=$1`)).
		WithArgs(pgUUID.String()).
		WillReturnError(errors.New("test error"))

	// exercise
	result, err := DeleteReservationData(context.Background(), mockConn, pgUUID.String())

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}

	if result != 0 {
		t.Fatal("expected 0 deleted rows, got", result)
	}

	if err = mockConn.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
