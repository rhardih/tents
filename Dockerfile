FROM golang:1.15-alpine

RUN mkdir /tents
COPY . /tents
WORKDIR /tents

RUN go build

CMD ["time", "/tents/tents"]
