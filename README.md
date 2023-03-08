# TODO

## Tracing
  - Рефактор кода

## Security

  - Вырезать password из access_token payload
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
