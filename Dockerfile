# Start with the official Golang image as a base
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download all the dependencies and build the application
RUN go get -d -v ./... && go build -o main .

# Expose the port on which your app will run
EXPOSE 8080

# Run the binary
CMD ["./main"]
