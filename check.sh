#!/usr/bin/env bash

set -euo pipefail

echo "Running gofmt check..."
UNFORMATTED="$(gofmt -l $(find . -type f -name '*.go'))"
if [[ -n "${UNFORMATTED}" ]]; then
  echo "The following files are not properly formatted:"
  echo "${UNFORMATTED}"
  exit 1
fi

echo "Running go vet (v5)..."
(cd "v5" && go vet ./...)

echo "Running go test (v5)..."
(cd "v5" && go test ./...)

echo "Running go test -race (v5)..."
(cd "v5" && go test -race ./...)

echo "Running go test (examples)..."
(cd "examples" && go test ./...)
