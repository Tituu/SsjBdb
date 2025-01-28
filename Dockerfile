# Use Golang as a base image for building the Go app
FROM golang:1.20-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go application into the container
COPY . .

# Build the Go app
RUN go build -o noti .

# Start a new stage to create a minimal runtime image
FROM alpine:latest

# Install necessary dependencies for running the Go binary (if any)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/noti .

# Expose the port your application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./noti"]
