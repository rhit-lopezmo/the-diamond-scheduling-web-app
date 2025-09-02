package db_utils

import (
	"context"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

func LoadTunnelData(ctx context.Context, conn IDBConn) ([]models.Tunnel, error) {
	var tunnels []models.Tunnel

	err := pgxscan.Select(
		ctx,
		conn,
		&tunnels,
		`SELECT * FROM tunnels`,
	)

	if err != nil {
		log.Println("[API] Error querying database:", err)
		return nil, err
	}

	return tunnels, nil
}
