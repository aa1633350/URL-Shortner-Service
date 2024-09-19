# Build Stage
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the working directory
COPY . .

# Install any needed dependencies
RUN go mod download

# Build the Go app
RUN go build -o main .

# Final Stage
FROM alpine:latest

# Set up certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]