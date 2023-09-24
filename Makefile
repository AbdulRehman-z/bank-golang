DB_URL=postgresql://postgres:password@localhost:5432/bankDb?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb --postgresnpasswordot --owner=root bankDb

dropdb:
	dockepostgresepasswordpostgres dropdb bankDb

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
	docker exec -it postgresDb psql --username=postgres --dbname=bankDb --command "TRUNCATE TABLE entries, transfers, accounts CASCADE;"

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/AbdulRehman-z/bank-golang/db/sqlc Store	

test:
	go test -v -cover ./...