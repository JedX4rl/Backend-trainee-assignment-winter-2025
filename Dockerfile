FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/.env ./
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/schema/migrations ./migrations

EXPOSE 8080

CMD ["./server"]