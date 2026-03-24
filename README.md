ORDER SERVICE
📌Описание
gRPC сервис для управления заказами
Возможности:

- CreateOrder
- GetOrder
- UpdateOrder
- DeleteOrder
- ListOrders
  Технологии:
- gRPC
- .env конфигурацию

⚙️ Конфигурация
Создайте .env файл с конфигурациями:
GRPC_PORT=50051
Network=tcp
Или воспользуйтесь env.example

🚀 Запуск

1. Установить зависимости
   go mod tidy
2. Запустить сервер
   go run cmd/server/main.go или make run

🧪 Генерация gRPC кода
make proto

🛠 Переменные окружения
Переменная Описание
GRPC_PORT Порт gRPC сервера
HTTP_PORT Порт http сервера
NETWORK Тип сети (tcp)

📦 Структура проекта
server/main.go
pkg/api/test (сгенерированный код)
