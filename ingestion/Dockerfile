FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ingestionbin ./ingestion/cmd

EXPOSE 8080
CMD ["./ingestionbin"]

