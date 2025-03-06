FROM golang:1.24.1-alpine3.21 as builder

WORKDIR /src

COPY go.mod /src/
COPY go.sum /src/
COPY cmd/ /src/cmd/

RUN go mod download \
    && GOOS=linux go build -v -o bin/sms github.com/itsalex/sentry-mattermost-sidecar/cmd/sms

FROM alpine:3.21

COPY --from=builder /src/bin/sms /usr/bin/go-sms

ENV GIN_MODE=release
ENV SMS_PORT=1323

EXPOSE 1323

CMD ["/usr/bin/go-sms"]