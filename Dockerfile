FROM golang:1.20.4-alpine3.17 as builder

WORKDIR /src

COPY go.mod /src/
COPY go.sum /src/
COPY cmd/ /src/cmd/

RUN go mod download \
    && GOOS=linux go build -v -o bin/sms github.com/itsalex/sentry-mattermost-sidecar/cmd/sms

FROM alpine:3.14.2

COPY --from=builder /src/bin/sms /usr/bin/go-sms

ENV GIN_MODE=release
ENV SMS_PORT=1323

EXPOSE 1323

CMD ["/usr/bin/go-sms"]