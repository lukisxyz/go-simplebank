PSQL_DOCKER := postgres14
PSQL_DBNAME := simplebank

create-postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=wap12345 -d postgres:14-alpine

postgres:
	docker start postgres14

createdb:
	docker exec -it $(PSQL_DOCKER) createdb --username=root --owner=root $(PSQL_DBNAME)

dropdb:
	docker exec -it $(PSQL_DOCKER) dropdb $(PSQL_DBNAME)

migrateup:
	migrate -path db/migration -database "postgresql://root:wap12345@127.0.0.1:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:wap12345@127.0.0.1:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: create-postgres postgres createdb dropdb migrateup migratedown