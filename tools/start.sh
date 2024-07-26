#! /bin/sh

echo "start wait-for-it"

set -e
echo "ensure connection"
while ! nc -z "${POSTGRES_HOST}" "${POSTGRES_PORT}"; do sleep 1; done

echo "start migrate"
migrate -path ./internal/migration -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" --verbose up

exec "$@"
