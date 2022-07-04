postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root user_api

dropdb:
	docker exec -it postgres12 dropdb user_api

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/user_api?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/user_api?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination=db/mock/querier.go -source=./db/sqlc/querier.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock