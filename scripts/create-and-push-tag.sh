#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/../"

BRANCH="$(git rev-parse --abbrev-ref HEAD)"

if [[ "$BRANCH" != "main" ]]; then
  echo 'Must be on main to create and push new tag!'
  exit 1
fi
if [[ -z $1 ]]; then
  echo "No version specified!"
  exit 1
fi

echo "Creating a new tag \"v${1}\""

git tag -a "v${1}" -m "v${1}"
git push origin "v${1}"
