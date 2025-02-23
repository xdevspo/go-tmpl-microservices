#!/bin/bash

# Путь к директории validate относительно корня проекта
VALIDATE_DIR="src/auth-service/vendor.protogen/validate"

# Проверяем существование директории
if [ ! -d "$VALIDATE_DIR" ]; then
    echo "Creating validate directory and downloading proto files..."

    # Создаем директорию с правильными правами
    mkdir -p "$VALIDATE_DIR"

    # Клонируем репозиторий во временную директорию
    git clone https://github.com/envoyproxy/protoc-gen-validate tmp_validate

    # Копируем proto файлы с сохранением прав доступа
    cp -p tmp_validate/validate/*.proto "$VALIDATE_DIR/"

    # Устанавливаем права доступа для созданных файлов
    chmod 644 "$VALIDATE_DIR"/*.proto

    # Удаляем временную директорию
    rm -rf tmp_validate

    echo "Proto files successfully initialized in $VALIDATE_DIR"
else
    echo "Directory $VALIDATE_DIR already exists. Skipping initialization."
fi