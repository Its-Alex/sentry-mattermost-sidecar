##
## Build
##
FROM golang:1.19-buster AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN export CGO_ENABLED=0 && go build -o /sentry-mattermost-sidecar ./main.go

##
## Production
##
FROM golang:1.18-alpine

WORKDIR /app/

COPY --from=build /sentry-mattermost-sidecar /app/sentry-mattermost-sidecar

ENTRYPOINT ["/app/sentry-mattermost-sidecar"]