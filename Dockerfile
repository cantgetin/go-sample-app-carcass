FROM golang:alpine

COPY . /go/src/app

WORKDIR /go/src/app/cmd

RUN go build -o app main.go

EXPOSE 8080

CMD ["./app"]