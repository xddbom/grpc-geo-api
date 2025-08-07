#!/bin/bash

set -e

# golangci-lint
echo "Running golangci-lint..."
if ! command -v golangci-lint &> /dev/null; then
  echo "❌ golangci-lint is not installed."
  exit 1
fi

if ! golangci-lint run ./...; then
  echo "❌ Linter has found problems."
  exit 1
fi
echo "✅ The code is clear."

# gofumpt
echo "Running gofumpt..."
if ! command -v gofumpt &> /dev/null; then
  echo "❌ gofumpt is not installed."
  exit 1
fi

UNFORMATTED=$(gofumpt -l .)
if [ -n "$UNFORMATTED" ]; then
  echo "❌ Code not formatted with gofumpt:"
  echo "$UNFORMATTED"
  exit 1
fi
echo "✅ Code formatted with gofumpt."