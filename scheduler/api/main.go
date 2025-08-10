package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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

	// connect to database with a 5 second timeout window before it fails + cancels
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
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
		tunnels := loadTunnelData()

		c.JSON(http.StatusOK, tunnels)
	})

	ginEngine.GET("/api/reservations", func(c *gin.Context) {
		reservations := loadReservationData()

		c.JSON(http.StatusOK, reservations)
	})

	ginEngine.GET("/api/reservations/search", func(c *gin.Context) {
		fromTimeStr := c.Query("from")
		toTimeStr := c.Query("to")
		tunnelIdStr := c.Query("tunnel_id")

		reservations := loadReservationDataWithParams(fromTimeStr, toTimeStr, tunnelIdStr)

		c.JSON(http.StatusOK, reservations)
	})

	ginEngine.POST("/api/reservations", func(c *gin.Context) {
		var reservation models.Reservation

		if err := c.BindJSON(&reservation); err != nil {
			log.Println("[API] Error binding JSON on POST method at /api/reservations.", err)
			return
		}

		c.Header("Location", "/api/reservations/"+reservation.Id)
		c.JSON(http.StatusCreated, reservation)
	})
}

func loadTunnelData() []models.Tunnel {
	var tunnels []models.Tunnel
	tunnels = append(tunnels, models.Tunnel{Id: 1, Name: "Tunnel 1"})
	tunnels = append(tunnels, models.Tunnel{Id: 2, Name: "Tunnel 2"})
	tunnels = append(tunnels, models.Tunnel{Id: 3, Name: "Tunnel 3"})
	tunnels = append(tunnels, models.Tunnel{Id: 5})

	return tunnels
}

func loadReservationData() []models.Reservation {
	var reservations []models.Reservation
	reservations = append(
		reservations,
		models.Reservation{
			Id:         "test-id",
			TunnelId:   1,
			CustomerId: "test-customer",
			Title:      "Lopez - 60 mins",
			StartsAt:   time.Date(2025, 8, 9, 12, 0, 0, 0, time.UTC),
			EndsAt:     time.Date(2025, 8, 9, 13, 0, 0, 0, time.UTC),
			Notes:      "Bring helmet",
		},
	)

	return reservations
}

func loadReservationDataWithParams(fromTime, toTime, tunnelId string) []models.Reservation {
	// create empty slice so it doesn't repsond with nil
	reservations := make([]models.Reservation, 0)

	if tunnelId == "3" {
		reservations = append(
			reservations,
			models.Reservation{
				Id:         "when-tunnel-id",
				TunnelId:   3,
				CustomerId: "test-customer",
				Title:      "Lopez - 60 mins",
				StartsAt:   time.Date(2025, 8, 9, 12, 0, 0, 0, time.UTC),
				EndsAt:     time.Date(2025, 8, 9, 13, 0, 0, 0, time.UTC),
				Notes:      "Bring helmet",
			},
		)

		return reservations
	}

	fromTimeParsed, err := time.Parse(time.RFC3339, fromTime)
	if err != nil {
		log.Println("[API] Could not parse from time when loading reservation data with params.", err)
		return reservations
	}

	if fromTimeParsed.After(time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)) {
		reservations = append(
			reservations,
			models.Reservation{
				Id:         "when-after-aug-first",
				TunnelId:   1,
				CustomerId: "test-customer",
				Title:      "Lopez - 60 mins",
				StartsAt:   time.Date(2025, 8, 9, 12, 0, 0, 0, time.UTC),
				EndsAt:     time.Date(2025, 8, 9, 13, 0, 0, 0, time.UTC),
				Notes:      "Bring helmet",
			},
		)

		return reservations
	}

	log.Println("[API] Found no matching reservations, returning nothing...")
	return reservations
}
