postgres:
	@docker rm -f postgres || true
	@docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -v pgdata:/var/lib/postgresql/data -d postgres

createdb:
	@docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	@docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
   sqlc generate  
test:	
   go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
