#! /bin/sh

echo "start wait-for-it"

set -e
echo "ensure connection"
while ! nc -z "${PGHOST}" "${POSTGRES_PORT}"; do sleep 1; done 

echo "start migrate"
migrate -path ./migration -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$PGHOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" --verbose up

echo "start app"
todo
exec "$@"
