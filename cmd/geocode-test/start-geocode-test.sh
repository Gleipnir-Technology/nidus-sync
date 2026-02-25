#!/run/current-system/sw/bin/bash
export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && air
