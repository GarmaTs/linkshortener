FROM golang:1.18.6-alpine3.16

RUN mkdir /app
ADD ./cmd/api /app
ADD . /app

WORKDIR /app

RUN go mod download
RUN go build -o main .
CMD ["/app/main"]