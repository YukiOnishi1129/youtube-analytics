#!/bin/bash

# Generate Go types from OpenAPI spec

# Directory paths
SPEC_DIR="../../spec/ingestion/http"
OUTPUT_DIR="internal/driver/http/generated"

# Ensure output directory exists
mkdir -p $OUTPUT_DIR

# Generate types from OpenAPI
oapi-codegen -generate types \
  -package generated \
  -o $OUTPUT_DIR/types.gen.go \
  $SPEC_DIR/dist/@typespec/openapi3/openapi.yaml

# Generate server interfaces
oapi-codegen -generate gin \
  -package generated \
  -o $OUTPUT_DIR/server.gen.go \
  $SPEC_DIR/dist/@typespec/openapi3/openapi.yaml

echo "HTTP types generated successfully!"