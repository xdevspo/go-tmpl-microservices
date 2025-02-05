# Определяем аргумент ENV (по умолчанию dev)
ARG ENV=dev

# Stage 1: Builder (используется только в test/prod)
FROM golang:1.23.6 AS builder
WORKDIR /app

# Загружаем зависимости (оптимизированное кеширование)
COPY src/org-service/go.mod src/org-service/go.sum ./
RUN go mod download

# Копируем код и собираем бинарник
COPY src/org-service /app
RUN go build -o service .

# Stage 2: Runtime (используется в test/prod)
FROM alpine:latest AS runtime
WORKDIR /app

# Копируем скомпилированный бинарник из builder-стадии
COPY --from=builder /app/service /app/service

# Копируем миграции и скрипт ожидания БД
COPY src/org-service/migrations /app/migrations

# Копируем общий `wait-for.sh` из `docker/`
COPY docker/wait-for.sh /app/wait-for.sh
RUN chmod +x /app/wait-for.sh

# Запускаем сервис с миграциями
CMD ["/app/wait-for.sh", "postgres-${ENV}", "5432", "--", "/app/service", "migrate", "up", "--database", "DB_DSN"]

# Stage 3: Dev (используется только в dev, без сборки!)
FROM golang:1.23.6 AS dev
WORKDIR /app

# Устанавливаем необходимые системные зависимости
RUN apt update && apt install -y \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем protoc-gen-go и protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Устанавливаем пакет для миграций БД
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Пакет для .env файлов
RUN go install github.com/joho/godotenv/cmd/godotenv@v1.5.1

# Экспортируем PATH для доступа к go-бинарникам
ENV PATH=$PATH:/go/bin

# Проверяем наличие protoc и плагинов
RUN protoc --version && protoc-gen-go --version && protoc-gen-go-grpc --version

# Копируем миграции и скрипт ожидания БД
COPY src/org-service/migrations /app/migrations

# Копируем общий `wait-for.sh` из `docker/`
COPY docker/wait-for.sh /app/wait-for.sh
RUN chmod +x /app/wait-for.sh

# Запускаем сервис с миграциями
CMD ["/app/wait-for.sh", "postgres-${ENV}", "5432", "--", "/app/service", "migrate", "up", "--database", "DB_DSN"]

# Устанавливаем `air` через `go install`
RUN go install github.com/air-verse/air@v1.61.7 && \
    mv /go/bin/air /usr/local/bin/air && \
    chmod +x /usr/local/bin/air

# Копируем общий `.air.toml` из `docker/`
COPY docker/.air.toml /app/.air.toml

# В dev код **не копируем** — он монтируется из локальной машины
CMD ["air", "-c", ".air.toml"]