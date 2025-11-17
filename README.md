## grpc-geo-api

A high-performance gRPC geo-data aggregation service. It fetches coordinates, weather, and POIs from various third-party providers. The project is implemented as a modular monolith, prepared for subsequent decomposition into microservices.

### Development Status: In progress 

## API (wip)

Below is a temporary mock response used during early development. 

#### **Health check**
```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health.Check
```
```json
{
  "status": "SERVING"
}
```

---

#### **Aggregate** (mock)
```bash
grpcurl -plaintext localhost:50051 aggregate.v1.Aggregator.Aggregate
```
```json
{
  "geo": "Kyiv, Ukraine",
  "weather": "Sunny 20°C"
}
```


## Quick start
#### **Installation**
```bash
git clone https://github.com/xddbom/grpc-geo-api
cd grpc-geo-api
```
#### **Dependencies**
```go
go mod tidy
```
#### **Run**
```go
go run ./services/edge-gateway/cmd  # Local

docker compose up -d --build        # Docker
```
#### **Makefiles**
```bash
make proto         
make compose-up    
make compose-down  
make run     
make docs
```
