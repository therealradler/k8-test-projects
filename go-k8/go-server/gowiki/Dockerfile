# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR  /app

COPY go.mod ./

RUN go mod download

COPY *.go ./

COPY *.txt ./

COPY *.html ./

RUN go build -o /go-wiki 

EXPOSE 8081

CMD ["/go-wiki"]