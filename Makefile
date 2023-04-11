all: clean buildserver

run:
	go run main.go

buildserver:
	go build ./cmd/auth_server

clean:
	rm -rf auth_server
