# Modules caching
FROM golang:1.17-alpine3.14 as modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Builder
FROM golang:1.17-alpine3.14 as builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/main.go

# Final
FROM scratch

COPY --from=builder /app/app.env /
COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

CMD ["/app"]