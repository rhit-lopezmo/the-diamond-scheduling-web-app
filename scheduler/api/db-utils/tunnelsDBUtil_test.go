package db_utils

import (
	"context"
	"errors"
	"regexp"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func Test_LoadTunnelData(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	rows := pgxmock.NewRows([]string{"id", "name", "is_active"}).
		AddRow(1, "Tunnel 1", true).
		AddRow(2, "Tunnel 2", true)

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tunnels`)).WillReturnRows(rows)

	// exercise
	result, err := LoadTunnelData(context.Background(), mockConn)

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

func Test_LoadTunnelData_Error(t *testing.T) {
	// setup
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())

	mockConn.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tunnels`)).WillReturnError(errors.New("test error"))

	// exercise
	_, err := LoadTunnelData(context.Background(), mockConn)

	// verify
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
