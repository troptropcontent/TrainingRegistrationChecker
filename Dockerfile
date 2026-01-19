# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o /app/checker .

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/checker .

CMD ["./checker"]
