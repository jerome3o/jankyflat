FROM golang:buster

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY main.go .
COPY index.html .

RUN go build .

CMD ["./producer"]