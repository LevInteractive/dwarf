#!/bin/bash
go build

# Redis Credentials
export REDIS_DB="0"
export REDIS_SERVER="localhost:32768"
export REDIS_PASS=""

# Port for the public application to listen.
export APP_PORT=":8000"

# Port for the internal gRPC server to listen.
export GRPC_PORT=":8001"

# The starting point for the short code algorithm.
export CHAR_FLOOR="2"

# URL to redirect to if the link isn't found.
export NOTFOUND_REDIRECT_URL="https://github.com"

# Full URL for the public dwarf app to run.
# No trailing slash.
export APP_BASE_URL="http://localhost:8000"

./dwarf
