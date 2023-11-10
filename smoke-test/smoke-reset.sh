#!/usr/bin/env sh

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl.yaml"}

set -e

. ./smoke.common.sh
trap cleanup EXIT

deleteCluster
createCluster
echo "* Applying"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
echo "* Resetting"
../cfctl reset --config "${CFCTL_CONFIG}" --debug --force
echo "* Done, cleaning up"
