PSQL_DOCKER := postgres14
PSQL_DBNAME := simplebank

DB_URL := postgresql://root:wap12345@127.0.0.1:5432/simplebank?sslmode=disable

create-postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=wap12345 -d postgres:14-alpine

postgres:
	docker start postgres14

createdb:
	docker exec -it $(PSQL_DOCKER) createdb --username=root --owner=root $(PSQL_DBNAME)

dropdb:
	docker exec -it $(PSQL_DOCKER) dropdb $(PSQL_DBNAME)

createmigrate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database ${DB_URL} -verbose up

migratedown:
	migrate -path db/migration -database ${DB_URL} -verbose down

migrateup1:
	migrate -path db/migration -database ${DB_URL} -verbose up 1

migratedown1:
	migrate -path db/migration -database ${DB_URL} -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

dev:
	go run main.go

createmock:
	mockery --dir=db/sqlc --output=db/mock --name=Store --outpkg=mocks

.PHONY: createmigrate migratedown1 migrateup1 createmock dev create-postgres postgres createdb dropdb migrateup migratedown