# Start with the official Golang image as a base
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

RUN go get -d -v ./...

# Build the Go application
RUN go build -o main .

# Expose the port on which your app will run
EXPOSE 8080

# Run the binary
CMD ["./main"]
