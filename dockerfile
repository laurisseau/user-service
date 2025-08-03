# Stage 1: Build the Go binary
FROM golang:1.21 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies early
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# IMPORTANT: Also copy the shared config module
COPY ../config ../config

# Set environment variable so Go uses the local config path
ENV GO111MODULE=on

# Build the application
RUN go build -o user-service main.go

# Stage 2: Use a minimal runtime image
FROM alpine:latest

# Create working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/user-service .

# Expose app port (change if needed)
EXPOSE 8080

# Run the binary
CMD ["./user-service"]
