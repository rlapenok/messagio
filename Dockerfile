FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN  go build -o main ./cmd/main.go

FROM golang:1.22.2-alpine3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml .
COPY --from=builder /app/.env .
EXPOSE 7070
CMD ["./main"]
