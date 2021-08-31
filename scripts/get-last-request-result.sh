#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

curl --location --request GET 'localhost:5000/__last_request__'
