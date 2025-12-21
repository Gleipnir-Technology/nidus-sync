#!/run/current-system/sw/bin/bash
PSQL_DSN="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" bob/bobgen-psql
