# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Set the current working directory inside the container
WORKDIR /usr/src/app

# Install Air (a live-reloading tool for Go)
RUN go install github.com/air-verse/air@latest

# Copy the Go modules files
COPY go.mod go.sum ./

# Download and install Go modules (dependencies)
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Run the application
# CMD ["air", "./cmd/main.go", "-b", "0.0.0.0"]
