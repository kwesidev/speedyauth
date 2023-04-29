FROM golang:1.20-alpine AS build_base
RUN apk add --no-cache git
# Set the Current Working Directory inside the container
WORKDIR /tmp/auth_server
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 
# Build the Go app
RUN go build ./cmd/auth_server
# Start fresh from a smaller image
FROM alpine:3.14
RUN apk add --no-cache bash
COPY --from=build_base /tmp/auth_server/auth_server /app/auth_server
COPY --from=build_base /tmp/auth_server/static/email_templates/* /static/email_templates/
USER root:root
RUN cd /app
# This container exposes port 8080 to the outside world
EXPOSE 8080
# Run the binary program produced by `go install`
CMD ["/app/auth_server"]