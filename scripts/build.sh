#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

docker-compose exec -T workspace go build -v -o bin/sms github.com/itsalex/sentry-mattermost-sidecar/cmd/sms
