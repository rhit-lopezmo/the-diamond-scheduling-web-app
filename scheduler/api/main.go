package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn

func main() {
	// ONLY FOR LOCALHOST, NEED TO REMOVE !!!
	_ = godotenv.Load()

	log.SetFlags(log.Ltime)

	// get all the env variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[API] Error finding 'PORT' in env file.")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("[API] Error finding 'DATABASE_URL' in env file.")
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		log.Fatal("[API] Error finding 'CORS_ORIGIN' in env file.")
	}

	// connect to database with a 5 second timeout window before it fails + cancels
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Print("[API] Error connecting to the database.", err)
		return
	}
	defer conn.Close(context.Background())

	// ping the database
	if err := conn.Ping(ctx); err != nil {
		log.Print("[API] Error when pinging the database", err)
		return
	}
	
	// setup gin
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
					"db": "down",
				},
			)
		}
		
		// ping success
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"status": "healthy",
				"db": "up",
			},
		)
	})
}
