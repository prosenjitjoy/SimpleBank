include .env

create_container:
	podman run --name ${CONTAINER_NAME} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -p 5432:5432 -d postgres

create_database:
	podman exec -it ${CONTAINER_NAME} createdb --username=${DB_USER} ${DB_NAME}

delete_database:
	podman exec -it ${CONTAINER_NAME} dropdb --username=${DB_USER} ${DB_NAME}

open_database:
	podman exec -it ${CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NAME}

create_migration:
	migrate create -ext sql -dir database/migration -seq create_tables

migrate_up:
	migrate -database ${MIGRATE_URL} -path database/migration -verbose up

migrate_up_last:
	migrate -database ${MIGRATE_URL} -path database/migration -verbose up 1

migrate_down:
	migrate -database ${MIGRATE_URL} -path database/migration -verbose down

migrate_down_last:
	migrate -database ${MIGRATE_URL} -path database/migration -verbose down 1

sqlc_generate:
	sqlc generate

mock_generate:
	mockgen -package mockdb -destination database/mockdb/store.go main/database/db Store

run_test:
	go test -v -cover ./...

run_server:
	go run main.go

.PHONY: create_container create_database delete_database open_database create_migration migrate_up migrate_up_last migrate_down migrate_down_last sqlc_generate mock_generate run_test run_server