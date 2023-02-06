FROM golang:1.20.0-alpine3.16

WORKDIR /app

COPY . .

RUN go build -o app main.go

EXPOSE 8080

CMD ["./app"]