#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

docker-compose up -d workspace
docker-compose exec -T workspace go mod download
./scripts/build.sh
docker-compose up -d
