FROM golang:1.9.2 AS builder

WORKDIR /go/src/github.com/cerra-io/base-cluster-state/

COPY main.go .
COPY ./cmd ./cmd
COPY ./clean ./clean
COPY ./vendor ./vendor
COPY ./update ./update
COPY ./utils ./utils
COPY ./vacuum ./vacuum
COPY ./server ./server

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cluster-state .

FROM alpine:latest
MAINTAINER 	Ali Al-Shabibi <ali#cerra.io>

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /go/src/github.com/cerra-io/base-cluster-state/cluster-state .

CMD [ "./cluster-state", "server" ]
