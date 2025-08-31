package db_utils

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

func Test_LoadCoachesData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	u1 := uuid.New()
	u2 := uuid.New()
	pgU1 := pgtype.UUID{Bytes: [16]byte(u1), Valid: true}
	pgU2 := pgtype.UUID{Bytes: [16]byte(u2), Valid: true}

	rows := pgxmock.NewRows([]string{"id", "first_name", "last_name"}).
		AddRow(pgU1, "John", "Doe").
		AddRow(pgU2, "Jane", "Doe")

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM coaches`)).WillReturnRows(rows)

	// exercise
	result, err := LoadCoachesData(context.Background(), mockConn)

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

func Test_LoadCoachesData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM coaches`)).WillReturnError(errors.New("test error"))

	// exercise
	_, err := LoadCoachesData(context.Background(), mockConn)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func Test_InsertCoachData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	query := `
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
	u := uuid.New()
	pgU := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	testCoach := models.Coach{
		Id:          pgU,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       nil,
		Phone:       "1112223333",
		Specialties: []string{models.SpecialtyHitting},
	}

	rows := pgxmock.NewRows([]string{"id", "first_name", "last_name", "email", "phone", "specialties"}).
		AddRow(
			testCoach.Id,
			testCoach.FirstName,
			testCoach.LastName,
			testCoach.Email,
			testCoach.Phone,
			testCoach.Specialties,
		)

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
		testCoach.FirstName,
		testCoach.LastName,
		testCoach.Phone,
		testCoach.Email,
		testCoach.Specialties,
	).WillReturnRows(rows)

	// exercise
	result, err := InsertCoachData(context.Background(), mockConn, testCoach)

	// verify
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if result == nil {
		t.Fatal("expected result, got none")
	}

	if result.Id != testCoach.Id {
		t.Fatal("expected id ", testCoach.Id, "got", result.Id)
	}
}

func Test_InsertCoachData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	query := `
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
	u := uuid.New()
	pgU := pgtype.UUID{Bytes: [16]byte(u), Valid: true}

	testCoach := models.Coach{
		Id:          pgU,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       nil,
		Phone:       "1112223333",
		Specialties: []string{models.SpecialtyHitting},
	}

	mockConn.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(
		testCoach.FirstName,
		testCoach.LastName,
		testCoach.Phone,
		testCoach.Email,
		testCoach.Specialties,
	).WillReturnError(errors.New("test error"))

	// exercise
	_, err := InsertCoachData(context.Background(), mockConn, testCoach)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
