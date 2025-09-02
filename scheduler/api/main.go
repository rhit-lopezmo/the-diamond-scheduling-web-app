package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	dbUtils "github.com/rhit-lopezmo/the-diamond-scheduling-web-app/api/db-utils"
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
	ginEngine.GET("/healthcheck", healthcheck)
	
	ginEngine.GET("/api/tunnels", getTunnels)
	
	ginEngine.GET("/api/reservations", getReservations)
	
	ginEngine.POST("/api/reservations", createReservation)
	
	ginEngine.GET("/api/reservations/:id", getReservationById)
	
	ginEngine.PUT("/api/reservations/:id", updateReservationById)
	
	ginEngine.DELETE("/api/reservations/:id", deleteReservationById)
	
	ginEngine.GET("/api/reservations/search", searchReservations)
	
	ginEngine.GET("/api/coaches", getCoaches)
	
	ginEngine.POST("/api/coaches", createCoach)
	
	ginEngine.PUT("/api/coaches/:id", updateCoachById)
	
	ginEngine.DELETE("/api/coaches/:id", deleteCoachById)
}

func healthcheck(c *gin.Context) {
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
}

func getTunnels(c *gin.Context) {
	tunnels, err := dbUtils.LoadTunnelData(c.Request.Context(), conn)
	
	if err != nil {
		log.Println("[API] Error loading tunnel data:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.JSON(http.StatusOK, tunnels)
}

func getReservations(c *gin.Context) {
	reservations, err := dbUtils.LoadReservationData(c.Request.Context(), conn)
	
	if err != nil {
		log.Println("[API] Error loading reservation data:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.JSON(http.StatusOK, reservations)
}

func createReservation(c *gin.Context) {
	var reservation models.Reservation
	
	if err := c.BindJSON(&reservation); err != nil {
		log.Println("[API] Error binding JSON on POST method at /api/reservations.", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}
	
	// send the data to the db
	result, err := dbUtils.InsertReservationData(c.Request.Context(), conn, reservation)
	if err != nil {
		log.Println("[API] Error inserting reservation:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.Header("Location", "/api/reservations/"+result.Id.String())
	c.JSON(http.StatusCreated, *result)
}

func getReservationById(c *gin.Context) {
	id := c.Param("id")
	
	reservation, err := dbUtils.LoadReservationById(c.Request.Context(), conn, id)
	if err != nil {
		log.Println("[API] Error loading reservation:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	if reservation == nil {
		log.Println("[API] Could not find reservation with id:", id)
		c.Status(http.StatusNotFound)
		return
	}
	
	c.JSON(http.StatusOK, *reservation)
}

func updateReservationById(c *gin.Context) {
	id := c.Param("id")
	
	var reservationUpdates models.ReservationUpdates
	
	if err := c.BindJSON(&reservationUpdates); err != nil {
		log.Println("[API] Error binding JSON on PUT method at /api/reservations/"+id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}
	
	reservation, err := dbUtils.UpdateReservationData(c.Request.Context(), conn, id, reservationUpdates)
	if err != nil {
		log.Println("[API] Error updating reservation:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	if reservation == nil {
		log.Println("[API] Cannot update reservation because it does not exist with id:", id)
		c.Status(http.StatusNotFound)
		return
	}
	
	c.JSON(http.StatusOK, reservation)
}

func deleteReservationById(c *gin.Context) {
	id := c.Param("id")
	
	rowsAffected, err := dbUtils.DeleteReservationData(c.Request.Context(), conn, id)
	
	if err != nil {
		log.Println("[API] Error deleting reservation:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// didn't delete anything
	if rowsAffected < 1 {
		log.Println("[API] Could not find reservation to delete with id:", id)
		c.Status(http.StatusNotFound)
		return
	}
	
	log.Println("[API] Successfully deleted reservation with id:", id)
	c.Status(http.StatusNoContent)
}

func searchReservations(c *gin.Context) {
	fromTimeStr := c.Query("from")
	toTimeStr := c.Query("to")
	tunnelIdStr := c.Query("tunnel_id")
	
	reservations, err := dbUtils.LoadReservationDataWithParams(c.Request.Context(), conn, fromTimeStr, toTimeStr, tunnelIdStr)
	if err != nil {
		log.Println("[API] Error loading reservation data:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.JSON(http.StatusOK, reservations)
}

func getCoaches(c *gin.Context) {
	coaches, err := dbUtils.LoadCoachesData(c.Request.Context(), conn)
	if err != nil {
		log.Println("[API] Error loading coaches:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.JSON(http.StatusOK, coaches)
}

func createCoach(c *gin.Context) {
	var coach models.Coach
	
	if err := c.BindJSON(&coach); err != nil {
		log.Println("[API] Error binding JSON on POST method at /api/coaches.", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}
	
	// send the data to the db
	result, err := dbUtils.InsertCoachData(c.Request.Context(), conn, coach)
	if err != nil {
		log.Println("[API] Error inserting coach:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	c.Header("Location", "/api/reservations/"+coach.Id.String())
	c.JSON(http.StatusCreated, *result)
}

func updateCoachById(c *gin.Context) {
	id := c.Param("id")
	
	var coachUpdates models.CoachUpdates
	
	if err := c.BindJSON(&coachUpdates); err != nil {
		log.Println("[API] Error binding JSON on PUT method at /api/coaches/"+id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return
	}
	
	coach, err := dbUtils.UpdateCoachData(c.Request.Context(), conn, id, coachUpdates)
	if err != nil {
		log.Println("[API] Error updating coach:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	if coach == nil {
		log.Println("[API] Cannot update coach because it does not exist with id:", id)
		c.Status(http.StatusNotFound)
		return
	}
	
	c.JSON(http.StatusOK, coach)
}

func deleteCoachById(c *gin.Context) {
	id := c.Param("id")
	
	rowsAffected, err := dbUtils.DeleteCoachData(c.Request.Context(), conn, id)
	
	if err != nil {
		log.Println("[API] Error deleting coach:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	
	// didn't delete anything
	if rowsAffected < 1 {
		log.Println("[API] Could not find coach to delete with id:", id)
		c.Status(http.StatusNotFound)
		return
	}
	
	log.Println("[API] Successfully deleted coach with id:", id)
	c.Status(http.StatusNoContent)
}
