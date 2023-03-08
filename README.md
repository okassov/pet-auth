# TODO

## Tracing
  - Рефактор кода
  - Заменить запись трейсов с файла на опентелеметри + темпо + графана

## Security

  - Добавить Redis для JWT и написать логику refresh-а

# Docker

```
docker build -t okassov/pet-auth:v1 .
```

# Swagger Init

Install Doc - https://github.com/swaggo/swag#declarative-comments-format

```
swag init -g internal/controller/http/v1/router.go
go run cmd/main.go
curl http://localhost:8080/swagger/index.html
```
