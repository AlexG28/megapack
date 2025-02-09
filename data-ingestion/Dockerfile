FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o data-ingestion ./cmd/main.go

EXPOSE 8080
EXPOSE 5432

CMD [ "./data-ingestion" ] 