FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/bot/main.go

FROM golang:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations

CMD ["./main"]
