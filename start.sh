#!/bin/sh

set -e 

file="/app/app.env"

echo "Checking for $file..."

if [ ! -f "$file" ]; then
  echo "Creating missing app.env file."
  touch "$file"
else
  echo "Found $file, sourcing it."
fi

# Load environment variables from app.env
set -a
. "$file"
set +a

echo "run db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"