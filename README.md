# Nidus Sync

This is the software that powers [Nidus Cloud Sync](https://sync.nidus.cloud).

## Building from source

First, you'll need [Nix](https://nix.dev).

Then:

```sh
nix develop
go build .
```

## Running

You'll need a number of environment variables for configuring things;

 * ARCGIS_CLIENT_ID - The client ID for ArcGIS oauth, configured with esri.
 * ARCGIS_CLIENT_SECRET - The client secret for ArcGIS oauth, configured with esri.
 * BASE_URL - The URL the site is hosted at, used for forming callback URLs. Should be complete like 'https://foo.bar'
 * BIND - The address and port to bind to. Use ':9001' for 'any address, port 9001'
 * ENVIRONMENT - either 'PRODUCTION' or 'DEVELOPMENT'. It's used to set things like oauth token length.
 * MAPBOX_TOKEN - The token to use with mapbox which is important for rendering maps.
 * POSTGRES_DSN - The DSN for connecting to the postgres database.
 * FIELDSEEKER_SCHEMA_DIRECTORY - The directory to write fieldseeker schema files for debugging.
 * USER_FILES_DIRECTORY - The directory for writing uploaded user data files (audio, images).

```sh
> BASE_URL=https://sync.nidus.cloud ARCGIS_CLIENT_ID=foo ARCGIS_CLIENT_SECRET=bar POSTGRES_DSN='postgresql://?host=/var/run/postgresql&dbname=nidus-sync' ./nidus-sync
```

### Districts

There's a table containing district information in the database, `import.district`. It was created with:

```
psql
CREATE SCHEMA import;
shp2pgsql -s 3857 -c -D -I CA_districts.shp import.district | psql -d nidus-sync
psql
ALTER TABLE import.district ADD COLUMN geom_4326 geometry(MultiPolygon,4326) GENERATED ALWAYS AS (ST_Transform(geom, 4326)) STORED;
```

## Hacking

### air

This project uses [air](https://github.com/air-verse/air) for fast compile-and-test loops. You can run it with:

```sh
> BASE_URL=https://sync.nidus.cloud ARCGIS_CLIENT_ID=foo ARCGIS_CLIENT_SECRET=bar POSTGRES_DSN='postgresql://?host=/var/run/postgresql&dbname=nidus-sync' air
```

### bob

This uses the bob query framework. You can regenerate the models for bob with:

```
PSQL_DSN="postgresql://dbname?host=/var/run/postgresql&sslmode=disable" go run github.com/stephenafamo/bob/gen/bobgen-psql@latest"
PSQL_DSN="postgresql://?host=/var/run/postgresql&sslmode=disable&dbname=nidus-sync" go run github.com/stephenafamo/bob/gen/bobgen-psql@latest
```

This will generate a bunch of files. They're already committed, you only need this if you change the database schema in some way.

### goose

This uses [goose](https://github.com/pressly/goose). You can use the goose command line to check status and make changes

```sh
> cd migrations
> GOOSE_DRIVER=postgres GOOSE_DBSTRING="dbname=nidus-sync sslmode=disable" goose status
> GOOSE_DRIVER=postgres GOOSE_DBSTRING="dbname=nidus-sync sslmode=disable" goose down
> GOOSE_DRIVER=postgres GOOSE_DBSTRING="dbname=nidus-sync sslmode=disable" goose up
```
