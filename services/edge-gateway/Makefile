IMAGE := edge-gateway:latest
CONTAINER := edge-gateway
COMPOSE_FILE := docker-compose.yml

.PHONY: build

build:
	docker build -t $(IMAGE) -f services/edge-gateway/Dockerfile .

run:
	docker run --rm --name $(CONTAINER) -p 50051:50051 -p 8081:8081 $(IMAGE)

stop:
	@docker rm -f $(CONTAINER) 2>/dev/null || true

rebuild: stop build run

logs:
	docker logs -f $(CONTAINER)

clean: stop
	@docker rmi $(IMAGE) 2>/dev/null || true

up:
	docker compose -f $(COMPOSE_FILE) up -d --build

down:
	docker compose -f $(COMPOSE_FILE) down

ps:
	docker ps --filter "name=$(CONTAINER)"

logs-compose:
	docker compose logs -f $(CONTAINER)

run-local:
	go run ./services/edge-gateway/cmd
