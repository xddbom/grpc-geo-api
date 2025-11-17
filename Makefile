SERVICE := edge-gateway
DOCKER_COMPOSE := docker compose

.PHONY: compose-up compose-down

compose-up:
	@echo "→ Starting external services..."
	$(DOCKER_COMPOSE) up -d --build

compose-down:
	@echo "→ Stopping external services..."
	$(DOCKER_COMPOSE) down

.PHONY: proto

proto:
	@echo "→ Generating protobuf code..."
	buf format
	buf lint
	buf generate

.PHONY: run

run:
	@echo "→ Running $(SERVICE)..."
	go run ./services/$(SERVICE)/cmd

.PHONY: $(SERVICES)

$(SERVICES):
	@echo "→ Building $@"
	docker build -t $@:latest -f services/$@/Dockerfile .
