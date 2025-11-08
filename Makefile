SERVICES := edge-gateway
.PHONY: $(SERVICES)

$(SERVICES):
	@echo "→ Building $@"
	docker build -t $@:latest -f services/$@/Dockerfile .
