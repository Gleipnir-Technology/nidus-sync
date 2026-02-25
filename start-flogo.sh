#!/run/current-system/sw/bin/bash
# normal
#export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && BIND=127.0.0.1:9001 FLOGO_BIND=:9000 FLOGO_UPSTREAM=http://127.0.0.1:9001 ../flogo/flogo -target .
# verbose
export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && BIND=127.0.0.1:9001 FLOGO_BIND=:9000 FLOGO_UPSTREAM=http://127.0.0.1:9001 FLOGO_VERBOSE=1 ../flogo/flogo -target .
# no TUI
#export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && BIND=127.0.0.1:9001 FLOGO_BIND=:9000 FLOGO_DISABLE_TUI=1 FLOGO_UPSTREAM=http://127.0.0.1:9001 ../flogo/flogo -target .
