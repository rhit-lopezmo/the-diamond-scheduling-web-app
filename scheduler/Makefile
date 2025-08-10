DEFAULT_GOAL: restart

.PHONY: up down restart

up:
	docker compose up --build -d

down:
	docker compose down -v

restart: down up
