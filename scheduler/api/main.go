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
					"db": "down",
				},
			)
			return
		}
		
		// ping success
		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "healthy",
				"db": "up",
			},
		)
	})
}
