#!/bin/bash

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl-dynamic.yaml"}

set -e

. ./smoke.common.sh
trap cleanup EXIT

deleteCluster
createCluster

echo "* Starting apply"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
echo "* Apply OK"

max_retry=5
counter=0
echo "* Verifying dynamic config reconciliation was a success"
until ../cfctl config status -o json --config "${CFCTL_CONFIG}" | grep -q "SuccessfulReconcile"; do
   [[ counter -eq $max_retry ]] && echo "Failed!" && exit 1
   echo "* Waiting for a couple of seconds to retry"
   sleep 5
   ((counter++))
done

echo "* OK"

echo "* Dynamic config reconciliation status:"
../cfctl config status --config "${CFCTL_CONFIG}"

echo "* Done"
