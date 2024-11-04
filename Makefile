postgres:
	docker run --name postgres14 --network tracking-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mahanta -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root tracking_inventory

dropdb: 
	docker exec -it postgres14 dropdb tracking_inventory

migrateup:
	migrate -path db/migration -database "postgresql://root:NlEJjxQQsAQzAawwDCpc9fWI0RXjNz0N@dpg-csjq3adsvqrc73ev90s0-a.oregon-postgres.render.com/tracking_inventory" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:mahanta@localhost:5432/tracking_inventory?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Yelsnik/trackinginventory/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server