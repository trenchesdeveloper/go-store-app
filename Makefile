postgres:
	docker run --name postgres14 -p 5438:5432 -e POSTGRES_PASSWORD=password -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=postgres go-market

dropdb:
	docker exec -it  postgres14 dropdb postgres12

migratecreate:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

migrateup:
	migrate -path internal/db/migration -database "postgresql://postgres:password@localhost:5438/go-market?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://postgres:password@localhost:5438/go-market?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	nodemon --watch './**/*.go' --signal SIGTERM --exec APP_ENV=dev 'go' run main.go