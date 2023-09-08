FROM golang:1.21.1

WORKDIR /app

COPY go.mod .
COPY central/main.go .

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]