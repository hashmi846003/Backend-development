#maketest:
 # echo  "make  test"
postgres:
	@docker rm -f postgres || true
	@docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -v pgdata:/var/lib/postgresql/data -d postgres

createdb:
	@docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	@docker exec -it postgres dropdb --username=root --owner=root simple_bank

.PHONY: postgres createdb dropdb

