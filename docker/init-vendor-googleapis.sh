#!/bin/bash

# Определяем директорию для хранения proto файлов Google API
GATEWAY_SERVICE_GOOGLE_PROTO_DIR="src/gateway-service/vendor.protogen/google"
TEMP_GOOGLEAPIS_DIR="tmp_googleapis"

# Функция для очистки временных файлов в случае ошибки
cleanup() {
    echo "Cleaning up temporary files..."
    rm -rf "$TEMP_GOOGLEAPIS_DIR"
}

# Устанавливаем обработчик ошибок
trap cleanup ERR

# Функция для проверки успешности выполнения команды
check_error() {
    if [ $? -ne 0 ]; then
        echo "Error: $1"
        cleanup
        exit 1
    fi
}

if [ ! -d "$GATEWAY_SERVICE_GOOGLE_PROTO_DIR" ]; then
    echo "Initializing Google API proto files..."

    # Создаем родительскую директорию, если она не существует
#    mkdir -p "vendor.protogen"
#    check_error "Failed to create vendor.protogen directory"

    # Клонируем репозиторий googleapis
    echo "Cloning googleapis repository..."
    git clone https://github.com/googleapis/googleapis "$TEMP_GOOGLEAPIS_DIR"
    check_error "Failed to clone googleapis repository"

    # Создаем целевую директорию для Google API
    mkdir -p "$GATEWAY_SERVICE_GOOGLE_PROTO_DIR"
    check_error "Failed to create google proto directory"

    # Перемещаем необходимые файлы
    echo "Moving API proto files..."
    cp -r "$TEMP_GOOGLEAPIS_DIR/google/api" "$GATEWAY_SERVICE_GOOGLE_PROTO_DIR/"
    check_error "Failed to move api directory"

    # Устанавливаем права доступа для файлов
    chmod -R 644 "$GATEWAY_SERVICE_GOOGLE_PROTO_DIR/api/"*.proto
    find "$GATEWAY_SERVICE_GOOGLE_PROTO_DIR" -type d -exec chmod 755 {} \;
    check_error "Failed to set file permissions"

    # Очищаем временные файлы
    echo "Cleaning up..."
    rm -rf "$TEMP_GOOGLEAPIS_DIR"
    check_error "Failed to remove temporary directory"

    echo "Google API proto files successfully initialized in $GATEWAY_SERVICE_GOOGLE_PROTO_DIR"
else
    echo "Directory $GATEWAY_SERVICE_GOOGLE_PROTO_DIR already exists. Skipping initialization."
fi