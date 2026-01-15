# syntax=docker/dockerfile:1

FROM golang:1.25

ENV CGO_ENABLED=0 

WORKDIR /build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -v -o /gatnbot-bin

RUN rm -rf /build

WORKDIR /gatnbot

CMD ["/gatnbot-bin"]
