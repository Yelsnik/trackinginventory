#!/bin/sh

set -e 

file="/app/app.env"

echo "run db migration"

if ! [ -f /app/app.env ]; then
  touch $file
  echo "file "$file" exists!"
  source /app/app.env
fi


source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"