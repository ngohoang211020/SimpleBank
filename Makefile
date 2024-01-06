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
	mockgen -package mockdb -destination db/mock/store.go github.com/ngohoang211020/simplebank/db/sqlc Store
gen-go:
	rm -f pb/**/*.go
	rm -f doc/swagger/*.swagger.json
	rm -f doc/statik/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
	proto/**/*.proto
	statik -src=./doc/swagger -dest=./doc -ns=simple_bank #Use to embed static file into golang code
evans:
	evans --host localhost --port 7070 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
asynqmon:
	docker run --name asynqmon -p 8079:8080 -d hibiken/asynqmon:0.6.1
.PHONY:	postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 new_migration db_docs db_schema mock sqlc test server gen-go evans redis asynqmon