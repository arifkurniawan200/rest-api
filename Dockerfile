# Start from the official Go image
FROM golang:1.19-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o main .

# Build a second stage image
FROM alpine:latest

# Install PostgreSQL client for database migration
RUN apk --no-cache add postgresql-client

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory inside the container
WORKDIR /app

# Run database migration
CMD ["sh", "-c", "go run main.go db:migrate up && ./main api"]

# Expose port 8080 for the API server
EXPOSE 8080
