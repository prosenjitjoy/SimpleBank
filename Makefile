include .env

create_container:
	podman run --name ${DB_CONTAINER} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -p 5432:5432 -d postgres

create_database:
	podman exec -it ${DB_CONTAINER} createdb --username=${DB_USER} ${DB_NAME}

delete_database:
	podman exec -it ${DB_CONTAINER} dropdb --username=${DB_USER} ${DB_NAME}

open_database:
	podman exec -it ${DB_CONTAINER} psql -U ${DB_USER} -d ${DB_NAME}

create_migration:
	migrate create -ext sql -dir database/migration -seq create_tables

migrate_up:
	migrate -database ${DATABASE_URL} -path database/migration -verbose up

migrate_up_last:
	migrate -database ${DATABASE_URL} -path database/migration -verbose up 1

migrate_down:
	migrate -database ${DATABASE_URL} -path database/migration -verbose down

migrate_down_last:
	migrate -database ${DATABASE_URL} -path database/migration -verbose down 1

sqlc_generate:
	sqlc generate

mock_generate:
	mockgen -package mockdb -destination database/mockdb/store.go main/database/db Store

proto_gererate:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank proto/*.proto

run_test:
	go test -v -cover ./...

run_server:
	go run main.go

dev_deploy:
	podman pod rm -af
	podman rm -af
	podman pod create -p 3000:3000 ${POD_NAME}
	podman pod start ${POD_NAME}
	podman run --pod ${POD_NAME} --name ${DB_CONTAINER} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASS} -e POSTGRES_DB=${DB_NAME} -d postgres
	podman build -t ${BE_CONTAINER}:latest .
	podman run --pod ${POD_NAME} --name ${BE_CONTAINER} -e DB_SOURCE=${MIGRATE_URL} ${BE_CONTAINER}:latest

.PHONY: create_container create_database delete_database open_database create_migration migrate_up migrate_up_last migrate_down migrate_down_last sqlc_generate mock_generate proto_gererate run_test run_server dev_deploy