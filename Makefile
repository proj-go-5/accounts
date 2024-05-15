DB_DSN ?= postgres://accounts:accounts@localhost:5432/accounts?sslmode=disable

migration_create:
	@if [ -z "$(MIGRATION_NAME)" ]; then \
		echo "Please provide a migration name using MIGRATION_NAME=<name>"; \
	else \
		migrate create -ext sql -dir internal/db/migrations/ -seq $(MIGRATION_NAME); \
	fi


migration_up:
	migrate -path internal/db/migrations/ -database "$(DB_DSN)" -verbose up

migration_down:
	migrate -path internal/db/migrations/ -database "$(DB_DSN)" -verbose down

create_admin:
	go run cmd/scripts/createadmin.go

start:
	go run cmd/server/main.go
