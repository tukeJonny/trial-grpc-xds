#!/usr/bin/env bash

set -euo pipefail

cd pb
protoc -I. --go_out=plugins=grpc:. ping.proto
