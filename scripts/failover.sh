#!/usr/bin/env bash
set -euo pipefail
echo "🚨 Killing echo1 to test breaker fail-over"
docker compose -f infra/docker-compose.yml kill -s SIGKILL echo1
sleep 3    # wait for breaker EWMA to trip