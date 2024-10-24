#!/bin/sh

set -e 

echo "run db migration"

# Load environment variables
. /app/app.env || { echo "Failed to load /app/app.env"; exit 1; }

# Check if DB_SOURCE is populated
if [ -z "$DB_SOURCE" ]; then
  echo "Error: DB_SOURCE is not set or empty"
  exit 1
else
  echo "DB_SOURCE is: $DB_SOURCE"
fi

/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"