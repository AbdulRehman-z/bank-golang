DB_URL=postgresql://postgres:password@localhost:5432/bankDb?sslmode=disable
CONTAINER_NAME=postgresDb
DB_NAME=bankDb

network:
	docker network create bank-network

postgres:
	docker run --name $(CONTAINER_NAME) --network bank-network -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it $(CONTAINER_NAME) createdb --username=postgres --owner=postgres $(DB_NAME)

dropdb:
	docker exec -it $(CONTAINER_NAME) dropdb $(DB_NAME)

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
	docker exec -it $(CONTAINER_NAME) psql --username=postgres --dbname=$(DB_NAME) --command "TRUNCATE TABLE entries, users, exchange_rates, verify_emails, transfers, sessions, accounts CASCADE;"

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/AbdulRehman-z/bank-golang/db/sqlc Store	

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc --openapiv2_opt=allow_merge=true,merge_file_name=bank \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

test:
	go test -v -cover ./...

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.2-alpine	

.PHONY: migrate-down migrate-up	