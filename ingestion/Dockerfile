FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ingestionbin ./ingestion/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ingestionbin .

EXPOSE 8080

CMD ["./ingestionbin"]

