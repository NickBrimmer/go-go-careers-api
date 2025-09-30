.PHONY: build up down restart logs clean dev

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	rm -rf tmp/

dev:
	docker-compose up --build

dev-logs:
	docker-compose up --build

test:
	go test ./...

.DEFAULT_GOAL := dev
