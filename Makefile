include .env
all: clean buildserver

run:
	go run main.go

buildserver:
	go build ./cmd/auth_server

test:
	go test ./...

dockerbuild:
	docker build -t authserver-app:latest .

migrateup:
	migrate -source file://db/migrations -database 'postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable' up

migratedown:
	migrate -source file://db/migrations -database 'postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable' down

clean:
	rm -rf auth_server
