FROM golang:1.15-alpine

WORKDIR /tents
COPY . .

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
