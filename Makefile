MIGRATE_BIN := ~/go/bin/migrate  # Adjust this to the actual path of your migrate binary
DATABASE_URL := postgresql://postgres:root@localhost:5431/personal_ai?query
MIGRATIONS_PATH := migrations
GO_BIN := /usr/local/go/bin/go

.PHONY: run
run:
	@$(GO_BIN) run .


.PHONY: create-migration
create-migration:
	@$(MIGRATE_BIN) create -ext sql -dir migrations $(name)


.PHONY: migrate-up 
migrate-up:
	@echo "Starting database migration up..."
	@$(GO_BIN) run . migrate-up
	@echo "Migration up completed successfully."


.PHONY: migrate-down 
migrate-down:
	@echo "Starting database migration down..."
	@$(GO_BIN) run . migrate-down
	@echo "Migration down completed successfully."
