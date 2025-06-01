# Use Go base image
FROM golang:1.24.3

# Set working directory
WORKDIR /usr/src/app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download






# Copy entire project
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the app
CMD ["./main"]
