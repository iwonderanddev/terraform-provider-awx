#!/usr/bin/env sh
set -eu

SCHEMA_URL="https://s3.amazonaws.com/awx-public-ci-files/awx/devel/schema.json"
OUTPUT_PATH="$(dirname "$0")/schema.json"

curl -fsSL "$SCHEMA_URL" -o "$OUTPUT_PATH"
echo "Updated $OUTPUT_PATH"
shasum -a 256 "$OUTPUT_PATH"
