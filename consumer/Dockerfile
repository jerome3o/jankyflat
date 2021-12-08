FROM golang:buster

RUN apt-get update && apt-get install -y libasound2-dev pkg-config

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY main.go .
COPY woo.mp3 .

RUN go build .

CMD ["./consumer"]