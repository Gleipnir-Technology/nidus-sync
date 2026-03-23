#!/run/current-system/sw/bin/bash
# MITM proxy
#MITM_PROXY=http://127.0.0.1:8080
# Flogo verbose
#FLOGO_VERBOSE=1
# No flogo TUI
#FLOGO_DISABLE_TUI=1

export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && 
	export $(cat .env | xargs) && \
	../flogo/flogo -target .
