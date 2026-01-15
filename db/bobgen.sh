#!/run/current-system/sw/bin/bash
PSQL_DSN="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" /tmp/bobgen-psql
#PSQL_DSN="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" bob/gen/bobgen-psql/bobgen-psql
