FROM golang:1.21 AS builder
WORKDIR /buildarea
COPY go.mod go.sum ./
# Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN mkdir /output
RUN CGO_ENABLED=0 GOOS=linux go build -o /output/app .


FROM ubuntu:latest AS production
# required to work with tls and aws
RUN apt-get update && apt-get install -y ca-certificates

# Set working directory to nginx asset directory
WORKDIR /app
COPY --from=builder /output/app .

# Expose port 8000 to the outside world
EXPOSE 8000

# Run the binary.
CMD ["./app"]
