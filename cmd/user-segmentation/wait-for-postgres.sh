#! /bin/sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h $host -d $POSTGRES_DB -U $POSTGRES_USER -c '\q'; do
  >&2 echo -e "Postgres is unavailable: sleeping...\n"
  sleep 1
done

>&2 echo -e "Postgres is up: executing $cmd\n"
exec $cmd