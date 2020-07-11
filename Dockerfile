FROM golang:1.14

WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon

COPY . /app
