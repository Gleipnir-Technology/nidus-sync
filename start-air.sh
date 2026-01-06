#!/run/current-system/sw/bin/bash
ARCGIS_CLIENT_ID=" " \
ARCGIS_CLIENT_SECRET=" " \
BASE_URL=" " \
BIND="127.0.0.1:9000" \
ENVIRONMENT="DEVELOPMENT" \
MAPBOX_TOKEN=" " \
POSTGRES_DSN="postgresql://?host=/var/run/postgresql&dbname=nidus-sync" \
FIELDSEEKER_SCHEMA_DIRECTORY=" " \
USER_FILES_DIRECTORY=" " \
export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && air
