.PHONY: build

default: build

build:
	GOOS=linux CGO_ENABLED=0 go build -o service

restart:
	docker-compose up -d --no-deps --build service
