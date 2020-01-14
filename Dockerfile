# Start from the latest golang base image
FROM golang:latest

RUN apt-get update
RUN apt-get upgrade -y

ENV GOBIN /go/bin

# Set the Current Working Directory inside the container
WORKDIR /app

RUN go get -u github.com/callicoder/go-docker

RUN go get -u github.com/nats-io/nats.go

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
