#!/run/current-system/sw/bin/bash
export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && ./nidus-sync 2>&1 | tee nidus-sync.log
