# Tools

Useful for doing one-off developer types of work. Can be run with:

```
psql -d nidus-sync -v org_id=3 -f tools/delete-org.sql
```

## Parcel data

You'll need to use `arcgis-go/example/layer-to-csv` in order to download the CSV files. Then copy them to a deployment server. Then run the import scripts.

```
dev$ cd arcgis-go/example/layer-to-csv
dev$ go build
dev$ ./layer-to-csv ...something I forget...
dev$ rsync -a *.csv {server}:/tmp
dev$ ssh {server}
server$ psql -d nidus-sync -f /tmp/tools/create-import-parcel-visalia.sql
server$ psql -d nidus-sync
nidus-sync=> \copy import.csv_parcel from '/tmp/parcel.csv' delimiters ',' csv header;
nidus-sync=> \q
server$ psql -d nidus-sync -f /tmp/tools/port-parcel-visalia.sql
server$ psql -d nidus-sync
nidus-sync=> GRANT SELECT ON parcel TO "tegola";
```

