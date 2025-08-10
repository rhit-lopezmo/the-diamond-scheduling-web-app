# Use this Makefile to manage the docker-compose stack

# Default target
.PHONY: up down restart

up:
	docker compose up --build -d

down:
	docker compose down

restart: down up
