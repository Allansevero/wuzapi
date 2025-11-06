FROM golang:1.24-alpine AS builder

# Install dependencies needed for building
RUN apk add --no-cache git ca-certificates

# Set the working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o wuzapi .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/wuzapi .

# Copy static files from builder stage
COPY --from=builder /app/static ./static

# Create necessary directories with proper permissions
RUN mkdir -p /root/dbdata && \
    mkdir -p /root/files && \
    chmod 777 /root/dbdata && \
    chmod 777 /root/files

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./wuzapi"]