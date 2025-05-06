# Makefile for Scheduling-system

APP_NAME=scheduling-system
DOCKER_IMAGE=$(APP_NAME):latest

.PHONY: build run test docker-build docker-run

## Build the Go application
build:
	go build -o bin/$(APP_NAME) ./...

## Run the Go application
run:
	go run cmd/main.go

## Run automated tests
test:
	go test ./... -v

## Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## Run Docker container
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE)
