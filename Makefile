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
	@if [ -z "$(LOGIN)" ]; then \
		echo "Please provide an admin login using LOGIN=<your_login>"; \
	elif [ -z "$(PASSWORD)" ]; then \
		echo "Please provide an admin password using PASSWORD=<your_password>"; \
	else \
		go run cmd/scripts/createadmin.go $(LOGIN) $(PASSWORD); \
	fi


start:
	go run cmd/server/main.go
