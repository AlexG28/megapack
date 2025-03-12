FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o monitoringbin ./monitoring/cmd/main.go

CMD [ "./monitoringbin" ] 