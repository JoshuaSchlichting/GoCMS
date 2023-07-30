# Start from the latest Golang base image
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o app .


FROM ubuntu:latest
# Set working directory to nginx asset directory
WORKDIR /app
COPY --from=builder /app/app /app/app
# Expose port 80
EXPOSE 80
# Run the binary.
CMD ["/app/app"]
