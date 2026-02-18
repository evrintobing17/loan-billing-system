# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o billing-api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/billing-api .

# Copy migrations folder (if needed for init scripts, though migrations are run by initdb in postgres container)
COPY migrations /migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./billing-api"]