#!/run/current-system/sw/bin/bash
export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && MITM_PROXY=http://127.0.0.1:8080 ./nidus-sync 2>&1 | tee nidus-sync.log
