postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

execdb:
	docker exec -it postgres-container psql -U root --dbname=simple_bank

createdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-container dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

test:
	go clean -testcache
	go test -v -cover ./...

.PHONY: postgres dropdb createdb migrateup migratedown test
