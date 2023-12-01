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
sqlc:
	sqlc	generate
test:
	go test -v -cover -short ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
.PHONY:	postgres migrateup migratedown mock