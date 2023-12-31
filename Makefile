postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

execdb:
	docker exec -it postgres-container psql -U root --dbname=simple_bank

createdb:
	docker exec -it postgres-container createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-container dropdb simple_bank

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

test:
	go clean -testcache
	go test -v -cover ./...

server:
	go run main.go

cleandb:
	docker exec postgres-container psql -U root --dbname=simple_bank \
	-c "DELETE FROM transfer WHERE 1=1;" \
	-c "DELETE FROM entry WHERE 1=1;" \
	-c "DELETE FROM account WHERE 1=1;" \
	-c "DELETE FROM users WHERE 1=1;"

mock:
	mockgen -package mockdb -destination db/mock/store.go tutorial.sqlc.dev/app/db/sqlc Store

sqlcgen:
	docker run --rm -v $(CURDIR):/src -w /src sqlc/sqlc generate 

.PHONY: postgres dropdb createdb migrateup migratedown migrateup1 migratedown1 test server cleandb mock