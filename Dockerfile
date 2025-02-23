# Use the official Go image for building
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN go build -o task-manager-api .

# Final container
FROM alpine:latest

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/task-manager-api .

# Set environment variables
ENV GO_ENV=production

# Specify the port
EXPOSE 8080

# Run the application
CMD ["./task-manager-api"]

