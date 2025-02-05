#!/bin/sh
# Ждём доступности сервиса перед запуском

# Разворачиваем переменную окружения, если передана в виде `postgres-${ENV}`
HOST=$(eval echo "$1")
PORT="$2"

echo "Waiting for $HOST:$PORT..."

while ! nc -z "$HOST" "$PORT"; do
  echo "Still waiting for $HOST:$PORT..."
  sleep 2
done

echo "$HOST:$PORT is available!"

# Убираем первые два аргумента (HOST и PORT)
shift 2

# Выполняем переданную команду
exec "$@"
