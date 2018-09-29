#!/bin/bash
go build

export CHAR_FLOOR=2
export APP_BASE_URL="http://localhost:8000"
./dwarf
