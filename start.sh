#!/bin/sh

set -e 

echo "run db migration $DB_SOURCE"
source /app/app.env
echo /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"