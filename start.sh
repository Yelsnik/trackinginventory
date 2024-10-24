#!/bin/sh

set -e 

AWS_REGION=${AWS_REGION:-"eu-north-1"}

echo "Fetching secrets from AWS Secrets Manager"
aws secretsmanager get-secret-value --region "$AWS_REGION" --secret-id tracking_inventory \
  --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > /app/app.env

# Verify that app.env is populated
if [ ! -s /app/app.env ]; then
  echo "Error: app.env is empty or missing"
  exit 1
fi

echo "Loaded environment variables:"
cat /app/app.env


echo "run db migration"
. /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up 

echo "start the app"
exec "$@"