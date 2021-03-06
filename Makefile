postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root reddit

dropdb:
	docker exec -it postgres12 dropdb reddit

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/reddit?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/reddit?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown
