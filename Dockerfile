FROM golang:1.17.0-alpine3.14 as builder

WORKDIR /src

COPY go.mod /src/
COPY go.sum /src/
COPY cmd/ /src/cmd/

RUN go mod download \
    && GOOS=linux GOARCH=amd64 go build -v -o bin/sms github.com/itsalex/sentry-mattermost-sidecar/cmd/sms

FROM alpine:3.14.2

COPY --from=builder /src/bin/sms /usr/bin/go-sms

ENV SMS_PORT=1323

EXPOSE 1323

CMD ["/usr/bin/go-sms"]