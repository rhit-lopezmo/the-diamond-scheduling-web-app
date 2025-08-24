package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/models"
)

var conn *pgx.Conn

func main() {
	_ = godotenv.Load()

	log.SetFlags(log.Ltime)

	// get all the env variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("[API] Error finding 'PORT' in env file.")
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		log.Fatalln("[API] Error finding 'CORS_ORIGIN' in env file.")
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		log.Fatalln("[API] Error finding 'POSTGRES_USER' in env file.")
	}

	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	if postgresPass == "" {
		log.Fatalln("[API] Error finding 'POSTGRES_PASSWORD' in env file.")
	}

	// generate dbUrl
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@db:5432/the-diamond-scheduler?sslmode=disable",
		postgresUser,
		postgresPass,
	)

	// setup + run migrations w/ goose
	sqlDB, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalln("[API] Error connecting to database:", err)
	}
	defer sqlDB.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		log.Println("[API] Error setting goose dialect:", err)
		return
	}

	if err = goose.Up(sqlDB, "db/migrations"); err != nil {
		log.Println("[API] Error running goose migrations:", err)
		return
	}

	log.Println("[API] Goose finished running migrations successfully.")

	// connect to database with a 5 second timeout window before it fails + cancels
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err = pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Println("[API] Error connecting to the database.", err)
		return
	}
	defer conn.Close(context.Background())

	// ping the database
	if err = conn.Ping(ctx); err != nil {
		log.Println("[API] Error when pinging the database", err)
		return
	}

	// setup gin
	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.Default()

	// setup CORS (for dev)
	ginEngine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// initial check by browser will go here before executing http method
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// setup the endpoints for the api
	setupEndpoints(ginEngine)

	// start server
	log.Printf("[API] Server started on port %s...\n", port)
	if err := ginEngine.Run(":" + port); err != nil {
		log.Println("[API] Error starting server.", err)
		return
	}

}

func setupEndpoints(ginEngine *gin.Engine) {
	ginEngine.GET("/healthcheck", func(c *gin.Context) {
		// ping DB to ensure it's up
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := conn.Ping(ctx); err != nil {
			// ping failed
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"status": "unhealthy",
					"db":     "down",
				},
			)
			return
		}

		// ping success
		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "healthy",
				"db":     "up",
			},
		)
	})

	ginEngine.GET("/api/tunnels", func(c *gin.Context) {
		tunnels, err := loadTunnelData(c.Request.Context())

		if err != nil {
			log.Println("[API] Error loading tunnel data:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, tunnels)
	})

	ginEngine.GET("/api/reservations", func(c *gin.Context) {
		reservations, err := loadReservationData(c.Request.Context())

		if err != nil {
			log.Println("[API] Error loading reservation data:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, reservations)
	})

	ginEngine.POST("/api/reservations", func(c *gin.Context) {
		var reservation models.Reservation

		if err := c.BindJSON(&reservation); err != nil {
			log.Println("[API] Error binding JSON on POST method at /api/reservations.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
			return
		}

		c.Header("Location", "/api/reservations/"+reservation.Id.String())
		c.JSON(http.StatusCreated, reservation)
	})

	ginEngine.GET("/api/reservations/:id", func(c *gin.Context) {
		id := c.Param("id")
		reservation := loadReservationById(id)

		c.JSON(http.StatusOK, reservation)
	})

	ginEngine.PUT("/api/reservations/:id", func(c *gin.Context) {
		id := c.Param("id")
		log.Println("[API] Temp log to use id for PUT request. id:", id)

		var reservation models.Reservation

		if err := c.BindJSON(&reservation); err != nil {
			log.Println("[API] Error binding JSON on PUT method at /api/reservations/:id.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, reservation)
	})

	ginEngine.DELETE("/api/reservations/:id", func(c *gin.Context) {
		id := c.Param("id")

		log.Println("[API] Temp log to delete reservation with id = ", id)

		c.Status(http.StatusNoContent)
	})

	ginEngine.GET("/api/reservations/search", func(c *gin.Context) {
		fromTimeStr := c.Query("from")
		toTimeStr := c.Query("to")
		tunnelIdStr := c.Query("tunnel_id")

		reservations := loadReservationDataWithParams(fromTimeStr, toTimeStr, tunnelIdStr)

		c.JSON(http.StatusOK, reservations)
	})
}

func loadTunnelData(ctx context.Context) ([]models.Tunnel, error) {
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

func loadReservationData(ctx context.Context) ([]models.Reservation, error) {
	var reservations []models.Reservation

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

func loadReservationDataWithParams(fromTime, toTime, tunnelId string) []models.Reservation {
	// create empty slice so it doesn't respond with nil
	reservations := make([]models.Reservation, 0)

	if tunnelId == "3" {
		reservations = append(
			reservations,
			models.Reservation{
				Id:        pgtype.UUID{},
				TunnelId:  &pgtype.UUID{},
				StartTime: pgtype.Timestamptz{},
				EndTime:   pgtype.Timestamptz{},
				Notes:     nil,
			},
		)

		return reservations
	}

	log.Println("[API] Found no matching reservations, returning nothing...")
	return reservations
}

func loadReservationById(id string) models.Reservation {
	reservation := models.Reservation{
		Id:       pgtype.UUID{},
		TunnelId: &pgtype.UUID{},
		Notes:    nil,
	}

	return reservation
}
