# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Arthur Coelho <arthur.cavalcante.puc@gmail.com>"

RUN mkdir /dummy_project

ADD . /dummy_project

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

WORKDIR /dummy_project

# Build the Go app
RUN go build -o main main.go

# Expose port 8000 to the outside world
EXPOSE 8888

# Command to run the executable
CMD ["/dummy_project/main"]