# Dockerfile
FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o main .

ENV APP_PORT="8080"

CMD ["/app/main", "-p", "8080"]
