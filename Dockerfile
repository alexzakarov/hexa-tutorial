# Start from the latest golang base image
FROM golang:1.21.3-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY config/config.prod.yaml ./
COPY docs/swagger.json ./
COPY i18n/tr.json ./
COPY i18n/en.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN swag init -g cmd/mono/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o service cmd/ticket/main.go

# Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# This lib is required
RUN apk add libc6-compat
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/service .
COPY --from=builder /app/config/config.prod.yaml ./config/config.prod.yaml
COPY --from=builder /app/docs/swagger.json ./doc/doc.json
COPY --from=builder /app/i18n/tr.json ./i18n/tr.json
COPY --from=builder /app/i18n/en.json ./i18n/en.json

ENV APP_ENV=prod

EXPOSE 8080:8080

# Command to run the executable
ENTRYPOINT ["./service"]
