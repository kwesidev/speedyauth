include .env
all: clean buildserver

run:
	go run main.go

buildserver:
	go build ./cmd/speedy_auth

test:
	go test ./...

dockerbuild:
	docker build -t speedyauth:latest .

migrateup:
	migrate -source file://db/migrations -database 'postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable' up

migratedown:
	migrate -source file://db/migrations -database 'postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable' down

clean:
	rm -rf speedy_auth
