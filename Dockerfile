FROM golang:alpine

RUN apk add musl-dev libc-dev gcc
RUN mkdir -p /go/src/wallet
ADD . /go/src/wallet
WORKDIR /go/src/wallet

RUN go get
RUN go build -o app .

# expose rest port
EXPOSE 3000

CMD ["/go/src/wallet/app"]
