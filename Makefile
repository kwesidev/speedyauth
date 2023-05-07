all: clean buildserver

run:
	go run main.go

buildserver:
	go build ./cmd/auth_server

dockerbuild:
	docker build -t authserver-app:latest .

clean:
	rm -rf auth_server
