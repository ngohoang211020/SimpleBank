DB_URL=postgresql://postgres:123456@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker-compose	up	-d
createdb:
	docker	exec	-it	golang-db-1	createdb	--username=postgres	simple_bank
dropdb:
	docker	exec	-it	golang-db-1	dropdb	simple_bank
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc:
	sqlc	generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

.PHONY:	postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration db_docs db_schema mock sqlc test server