#!/bin/bash
go build

export REDIS_DB="0"
export REDIS_SERVER="localhost:6379"
export REDIS_PASS=""
export APP_PORT=":8000"
export GRPC_PORT=":8001"
export CHAR_FLOOR="2"
export NOTFOUND_REDIRECT_URL="https://github.com"
export APP_BASE_URL="http://localhost:8000"
./dwarf
