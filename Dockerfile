FROM golang:1.14.4-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/akali/steplems-bot

COPY . .

ENV CGO_ENABLED=0

RUN go get -v all
WORKDIR app
RUN go build -o /bin/steplemsbot

FROM alpine

WORKDIR /bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/steplemsbot /bin/steplemsbot
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait
CMD ["/bin/steplemsbot"]
