# my-shop

Демонстрационный сервис для отображения данных о заказе (Go + PostgreSQL + Kafka).

## Быстрый старт

1. **Установите Docker и docker-compose**

2. **Запустите инфраструктуру:**

```sh
docker-compose up -d
```

3. **Инициализируйте базу данных:**

```sh
docker exec -i myshop_postgres psql -U myshop_user -d myshop_db < init.sql
```

4. **Установите Go (1.21+)**

5. **Скачайте зависимости:**

```sh
go mod tidy
```

6. **Запустите сервис:**

```sh
go run ./cmd/main.go
```

7. **Откройте в браузере:**

[http://localhost:8081/](http://localhost:8081/)

---

## Kafka

- Сервис автоматически слушает топик `orders` в Kafka (`localhost:9092`).
- Для теста можно отправить сообщение в Kafka с помощью любого клиента (например, [kcat](https://github.com/edenhill/kcat)) или скрипта.

## Структура проекта
- `cmd/main.go` — точка входа
- `internal/model` — модели данных
- `internal/db` — работа с БД
- `internal/cache` — кэш заказов
- `internal/kafka` — consumer Kafka
- `internal/server` — HTTP сервер и веб-интерфейс

## Пример сообщения заказа (JSON)
См. файл `model.json` в проекте.

---

**P.S.** Если что-то не работает — проверь, не включён ли VPN (для Docker).
