DB_URL := 'postgres://postgres:rahasia@localhost:5432/user_service?sslmode=disable'

.PHONY: migration
migration:
	migrate create -ext sql -dir database/migration -seq ${name}

.PHONY: migrate-up
migrate-up:
	migrate -database ${DB_URL} -path database/migration up

.PHONY: migrate-down
migrate-down:
	migrate -database ${DB_URL} -path database/migration down
