FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR $GOPATH/src/github.com/Syfaro/tg-webhook-prom/
COPY . .
RUN go get -d -v
RUN go build -o /go/bin/tg-webhook-prom

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/tg-webhook-prom /go/bin/tg-webhook-prom
CMD ["/go/bin/tg-webhook-prom"]
