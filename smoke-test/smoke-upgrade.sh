#!/usr/bin/env bash

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl.yaml"}

set -e

. ./smoke.common.sh
trap cleanup EXIT

deleteCluster
createCluster

remoteCommand() {
  local userhost="$1"
  shift
  echo "* Running command on ${userhost}: $*"
  bootloose ssh "${userhost}" -- "$*"
}

# Create config with older version and apply
K0S_VERSION="${K0S_FROM}"
echo "Installing ${K0S_VERSION}"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
remoteCommand "root@manager0" "k0s version | grep -q ${K0S_FROM}"

K0S_VERSION=$(curl -s "https://docs.k0sproject.io/stable.txt")

# Create config with latest version and apply as upgrade
echo "Upgrading to k0s ${K0S_VERSION}"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
remoteCommand "root@manager0" "k0s version | grep -q ${K0S_VERSION}"
