# Stage 1: Build the Go binary
FROM golang:1.23.5

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies early
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Enable Go modules (usually enabled by default)
ENV GO111MODULE=on

# Build the application binary
RUN go build -o user-service main.go

# Stage 2: Use a minimal runtime image
FROM debian:bookworm-slim

# Install CA certificates for HTTPS
RUN apt-get update && apt-get install -y ca-certificates \
    && update-ca-certificates \

# Set working directory
WORKDIR /root/

# Copy the compiled binary from builder stage
COPY --from=0 /app/user-service .

# Expose application port
EXPOSE 8082

# Run the binary
CMD ["./user-service"]
