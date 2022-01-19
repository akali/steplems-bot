FROM alpine:3 as certs
RUN apk --no-cache add ca-certificates

FROM golang:1.17-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/akali/steplems-bot

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0

RUN go get -v all

COPY . .

WORKDIR app
RUN go build -o /bin/steplemsbot

FROM alpine

WORKDIR /bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/steplemsbot /bin/steplemsbot
CMD ["/bin/steplemsbot"]
