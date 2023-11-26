FROM golang:1.21.4

WORKDIR /var/www/service-poll

COPY app/ ./

RUN go mod download

WORKDIR /var/www/service-poll/main

RUN go build -o ../service-poll

WORKDIR /var/www/service-poll

EXPOSE 5000

CMD ["./service-poll"]