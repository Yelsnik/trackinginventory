#!/bin/sh

set -e 

source /app/app.env
echo "run db migration $DB_SOURCE"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"