# Step 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules files to the container and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go source code to the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Step 2: Create a minimal image for running the application
FROM alpine:latest

# Install necessary dependencies (e.g., for glibc if needed)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application listens on
EXPOSE 8080

# Run the Go application
CMD ["./main"]
