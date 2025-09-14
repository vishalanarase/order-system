# Order Processing System with Kafka

A production-ready, event-driven order processing system built with Go and Apache Kafka. This system demonstrates modern microservices architecture patterns for handling e-commerce orders with scalability, reliability, and fault tolerance.

## 🚀 Features

- **RESTful HTTP API** - Gin-based API for order creation and health checks
- **Event-Driven Architecture** - Kafka-based messaging for loose coupling
- **Order Processing** - Async order validation and processing
- **Payment Processing** - Simulated payment gateway integration
- **Notification Service** - Email and SMS notifications for order events
- **Scalable Design** - Consumer groups and horizontal scaling support
- **Production Ready** - Structured logging, configuration management, graceful shutdown
- **Docker Support** - Local development with Kafka infrastructure

## 📋 System Architecture

```
HTTP Client → Order API → Kafka → [Order Processor → Payment Processor → Notification Service]
```

### Components

- **Order API**: REST endpoints for order creation (`POST /orders`)
- **Order Processor**: Handles order validation and business logic
- **Payment Processor**: Simulates payment processing (85% success rate)
- **Notification Service**: Sends email/SMS notifications for order events
- **Kafka**: Message broker with multiple topics for event streaming

### Event Flow

1. `ORDER_CREATED` → Order created via HTTP API
2. `ORDER_PROCESSING` → Order being processed
3. `ORDER_PAID`/`PAYMENT_FAILED` → Payment result
4. Notifications sent for each event type

## 🛠️ Technology Stack

- **Go 1.21+** - Backend programming language
- **Apache Kafka** - Message streaming platform
- **Gin** - HTTP web framework
- **Docker & Docker Compose** - Containerization and local development
- **Zap** - Structured logging
- **Viper** - Configuration management

## 📦 Installation

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/order-system.git
   cd order-system
   ```

2. **Start Kafka infrastructure**
   ```bash
   docker-compose up -d
   ```

3. **Build and run the application**
   ```bash
   go mod tidy
   go build -o order-system ./cmd
   ./order-system
   ```

4. **Verify services are running**
   ```bash
   curl http://localhost:8080/health
   ```

## 🚀 Usage

### Create an Order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "customer_123",
    "email": "customer@example.com",
    "items": [
      {
        "product_id": "prod_001",
        "product_name": "Premium Widget",
        "quantity": 2,
        "price": 29.99
      },
      {
        "product_id": "prod_002", 
        "product_name": "Standard Widget",
        "quantity": 1,
        "price": 15.99
      }
    ]
  }'
```

### Response
```json
{
  "message": "Order created successfully",
  "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status": "CREATED"
}
```

### Health Check

```bash
curl http://localhost:8080/health
```

## ⚙️ Configuration

The system uses a `config.yaml` file for configuration:

```yaml
kafka:
  brokers:
    - "localhost:9092"
  group_id: "order-service-group"
  topics:
    - "orders"
    - "payments"
    - "notifications"
  client_id: "order-service"

http:
  port: "8080"

notification:
  email_enabled: true
  sms_enabled: false
  email_from: "noreply@ordersystem.com"

service_name: "order-processing-system"
environment: "development"
```

### Environment Variables

All configuration values can be overridden with environment variables using the `SPF13_VIPER` convention:

```bash
export HTTP_PORT=8081
export KAFKA_BROKERS_0="kafka:9092"
export ENVIRONMENT="production"
```

## 🐳 Docker Development

### Start Kafka Cluster

```bash
docker-compose up -d
```

This starts:
- Zookeeper (port 2181)
- Kafka broker (port 9092)
- Kafka UI (port 8081)

### Access Kafka UI

Open http://localhost:8081 to monitor Kafka topics and messages.

### Stop Services

```bash
docker-compose down
```

## 📊 Monitoring and Logging

### Structured Logging

The system uses Zap for structured JSON logging:

```json
{
  "level": "info",
  "ts": "2023-12-07T10:30:00Z",
  "msg": "Order created and sent to Kafka",
  "order_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "customer_123",
  "total_amount": 75.97
}
```

### Log Levels

- `DEBUG`: Detailed processing information
- `INFO`: Service operations and events
- `WARN`: Potential issues
- `ERROR`: Operation failures

## 🔧 Development

### Project Structure

```
order-system/
├── cmd/                 # Main application entry point
├── internal/           # Private application code
│   ├── api/           # HTTP handlers and routes
│   ├── kafka/         # Kafka producer/consumer implementations
│   ├── models/        # Data structures
│   └── services/      # Business logic services
├── pkg/               # Public reusable packages
│   ├── config/        # Configuration loading
│   └── logger/        # Logging setup
├── config.yaml        # Application configuration
└── docker-compose.yml # Local development infrastructure
```

### Adding New Features

1. **New Event Types**: Add to `internal/models/order.go`
2. **New Services**: Create in `internal/services/`
3. **New API Endpoints**: Add to `internal/api/handlers.go`
4. **New Configuration**: Update `config.yaml` and `pkg/config/config.go`

### Testing

```bash
# Run tests
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🚀 Deployment

### Build for Production

```bash
GOOS=linux GOARCH=amd64 go build -o order-system ./cmd
```

### Docker Production Image

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o order-system ./cmd

FROM alpine:latest
COPY --from=builder /app/order-system /app/
COPY config.yaml /app/
EXPOSE 8080
CMD ["/app/order-system"]
```

### Kubernetes Deployment

Example deployment manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-system
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: order-system
        image: order-system:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: KAFKA_BROKERS_0
          value: "kafka-cluster:9092"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or issues:

1. Check the [Issues](https://github.com/your-username/order-system/issues) page
2. Create a new issue with detailed description
3. Contact: your-email@example.com

## 🙏 Acknowledgments

- Apache Kafka community
- Gin Web Framework
- Uber Zap logging library
- Viper configuration library

---

**Note**: This is a demonstration system. For production use, consider adding:
- Database persistence
- Authentication/Authorization
- Rate limiting
- Advanced monitoring (Prometheus, Grafana)
- SSL/TLS encryption
- Backup and recovery procedures