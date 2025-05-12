# Система обрабатывает 20k метрик в секунду по gRPC и отображает их через Prometheus

## Использование
```
# Сборка
bash build/build.sh
sudo docker compose -f ./docker/docker-compose.yml build
sudo docker compose -f ./docker/docker-compose.yml start

# Генератор метрик
./docker/sender/sender

# Отображение результата
xbrowser http://localhost:9090
```
🦝 🐭 🐹 🐰 🐿️ 🦦 🦡

## Дерево проекта

```
.
├── build
│   └── build.sh
├── cmd
│   ├── grpc-receiver
│   │   └── main.go
│   ├── metrics
│   │   └── main.go
│   └── sender
│       └── main.go
├── coverage.html
├── docker
│   ├── docker-compose.yml
│   ├── grpc-receiver
│   │   ├── Dockerfile
│   │   └── grpc-receiver
│   ├── metrics
│   │   ├── Dockerfile
│   │   └── metrics
│   ├── prometheus
│   │   ├── Dockerfile
│   │   └── prometheus.yml
│   ├── redis
│   │   ├── Dockerfile
│   │   └── redis.conf
│   └── sender
│       └── sender
├── gen
│   └── hamsterbeat.grpc
│       ├── hamsterbeat_grpc.pb.go
│       ├── hamsterbeat.pb.go
│       └── mocks_test.go
├── go.mod
├── go.sum
├── internal
│   └── hamsterbeat
│       ├── grpc.go
│       ├── grpc_test.go
│       ├── hamsterbeat.go
│       ├── hamsterbeat_test.go
│       ├── metrics.go
│       ├── metrics_test.go
│       └── mocks_test.go
├── LICENSE
├── proto
│   └── hamsterbeat.proto
└── README.md
```
