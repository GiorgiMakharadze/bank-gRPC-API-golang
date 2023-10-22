DB_URL=postgresql://root:secret@localhost:5432/bank?sslmode=disable

postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root bank
dropdb:
	docker exec -it postgres16 dropdb bank	
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up
migrateupone:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down	
migratedownone:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml	
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/GiorgiMakharadze/bank-API-golang/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown migrateupone migratedownone sqlc test server mock db_docs db_schema
	