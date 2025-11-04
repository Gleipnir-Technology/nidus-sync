# Nidus Sync

This is the software that powers [Nidus Cloud Sync](https://sync.nidus.cloud).

## Building from source

First, you'll need [Nix](https://nix.dev).

Then:

```sh
nix develop
go build .
```

## Hacking

### bob

This uses the bob query framework. You can regenerate the models for bob with:

```
PSQL_DSN="postgresql://dbname?host=/var/run/postgresql&sslmode=disable" go run github.com/stephenafamo/bob/gen/bobgen-psql@latest"
PSQL_DSN="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" go run github.com/stephenafamo/bob/gen/bobgen-psql@latest
```

This will generate a bunch of files. They're already committed, you only need this if you change the database schema in some way.
