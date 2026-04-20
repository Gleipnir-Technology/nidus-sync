#!/run/current-system/sw/bin/bash
# with MITM
# export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && MITM_PROXY=http://127.0.0.1:8080 ./nidus-sync 2>&1 | tee nidus-sync.log
#
# original recipe
#export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && ./nidus-sync 2>&1 | tee nidus-sync.log

# force production environment
# export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && ./nidus-sync -prod 2>&1 | tee nidus-sync.log
#
# force production environment, but with debug logging
 export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && \
	export $(cat .env | xargs) && \
	./nidus-sync -prod
#
# Use nix build output, force production environment
#export $(cat /var/run/secrets/nidus-dev-sync-env | xargs) && ./result/bin/nidus-sync -prod 2>&1 | tee nidus-sync.log
