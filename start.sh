#!/bin/sh

set -e

echo "running db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "starting the app"
exec "$@"