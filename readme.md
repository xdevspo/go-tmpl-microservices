## Шаблон проекта на микросервисах
  
Для всех микросервисов используется одна БД PostgreSQL.  

### Запуск с окружением
```bash
docker-compose --env-file .env.dev build --no-cache  
docker-compose --env-file .env.dev up
```

>Если работаете под Windows, то в настройках IDE указать что перенос строки LF.  

### Отключить docker-compose.override.yml 
```bash
docker-compose -f docker-compose.yml --env-file .env.test up
```  

### Генерация proto-файлов
В сервисе выполнить команду  
```bash
make generate-api
```
> По умолчанию генерация proto-файлов происходит по все папкам в папке api/.  
> Чтобы сгенерировать определенный api, при выполнении команды добавить параметр API={folder}, где  
> folder - папка с proto-файлами

По умолчанию будет генерация api версии 1.  
Для генерации другой версии  

```bash
make generate-api VERSION=v2
``` 

### БД

Создании миграций
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```
