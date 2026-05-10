FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /src

COPY go.mod /src/
COPY go.sum /src/
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/

RUN go mod download \
    && GOOS=linux go build -v -o bin/sms github.com/itsalex/sentry-mattermost-sidecar/cmd/sms

FROM alpine:3.22

COPY --from=builder /src/bin/sms /usr/bin/go-sms

ENV GIN_MODE=release
ENV SMS_PORT=1323

EXPOSE 1323

CMD ["/usr/bin/go-sms"]