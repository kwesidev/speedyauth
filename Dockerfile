FROM golang:1.20-alpine AS build_base
RUN apk add --no-cache git
# Set the Current Working Directory inside the container
WORKDIR /tmp/speedy_auth
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 
# Build the Go app
RUN go build ./cmd/speedy_auth
# Start fresh from a smaller image
FROM alpine:3.14
RUN apk add --no-cache bash
RUN apk add tzdata
COPY --from=build_base /tmp/speedy_auth/speedy_auth /app/speedy_auth
COPY --from=build_base /tmp/speedy_auth/static/email_templates/* /static/email_templates/
RUN cp /usr/share/zoneinfo/Africa/Johannesburg /etc/localtime
RUN echo "Africa/Johannesburg" >  /etc/timezone
RUN cd /app
# This container exposes port 8080 to the outside world
EXPOSE 8080
# Run the binary program produced by `go install`
CMD ["/app/speedy_auth"]