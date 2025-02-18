# We use the official GO image for assembly
FROM golang:1.23-alpine AS builder

# Install the working directory
WORKDIR /app

# Copy Go.mod and Go.sum files to install dependencies
COPY go.mod go.sum ./

# Set dependencies
RUN go mod download

# Copy the entire project to the container
COPY . .

# We collect the application
RUN go build -o task-manager-api .

# We use the minimum Alpine image for the final container
FROM alpine:latest

# Install the working directory
WORKDIR /app

# Copy the collected binary file
COPY --from=builder /app/task-manager-api .

# Copy the .env file
COPY .env . 

# Indicate the port that will use the application
EXPOSE 8080

# Team to launch the application
CMD ["./task-manager-api"]
