# Переменная версии (по умолчанию `v1`)
VERSION ?= v1
# Переменная окружения (по умолчанию `dev`)
ENV ?= dev
ENV_FILE = internal/config/.env.$(ENV)

# Базовые пути
PROTO_SRC=api/$(VERSION)
PROTO_OUT=pkg/$(VERSION)

# API, all - генерирует все файлы в папке PROTO_SRC (по умолчанию по всем папкам)
API ?= all

# Команда генерации gRPC API
generate-api:
	@mkdir -p $(PROTO_OUT)
	@if [ "$(API)" = "all" ]; then \
		proto_files=$$(find $(PROTO_SRC) -type f -name "*.proto"); \
	else \
		proto_files=$$(find $(PROTO_SRC)/$(API) -type f -name "*.proto"); \
	fi; \
	protoc --proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$$proto_files

# Миграции UP
migrate-up:
	set -a && export $$(grep -v '^#' $(ENV_FILE) | xargs) && set +a && \
    migrate -path ./migrations -database "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" up

# Миграции DOWN
migrate-down:
	set -a && export $$(grep -v '^#' $(ENV_FILE) | xargs) && set +a && \
    migrate -path ./migrations -database "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" down