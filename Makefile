DB_URL=postgresql://postgres:password@localhost:5432/bankDb?sslmode=disable
CONTAINER_NAME=postgresDb
DB_NAME=bankDb

network:
	docker network create bank-network

postgres:
	docker run --name $(CONTAINER_NAME) --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it $(CONTAINER_NAME) createdb --postgresnpasswordot --owner=root $(DB_NAME)

dropdb:
	dockepostgresepasswordpostgres dropdb $(DB_NAME)

migrate-up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

rmtestcache:
	go clean -testcache

truncate:
	docker exec -it $(CONTAINER_NAME) psql --username=postgres --dbname=$(DB_NAME) --command "TRUNCATE TABLE entries, users, exchange_rates, transfers, accounts CASCADE;"

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/AbdulRehman-z/bank-golang/db/sqlc Store	

test:
	go test -v -cover ./...

.PHONY: migrate-down migrate-up	