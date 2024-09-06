# Use the official Golang image as the base image
FROM golang:1.23.0-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./api ./api

# Build the Go app
RUN go build -o /go-shorty ./api/main.go

# Use a minimal base image
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the pre-built binary from the builder stage
COPY --from=builder /go-shorty /go-shorty

# Copy the .env file
COPY ./api/.env /app/.env

# Set permissions for the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Set the working directory
WORKDIR /app

# Expose the application port
EXPOSE 3000

# Command to run the executable
CMD ["/go-shorty"]