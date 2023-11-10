#!/usr/bin/env sh

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl-openssh.yaml"}

set -e

. ./smoke.common.sh
trap cleanup_openssh EXIT

cleanup_openssh() {
  cleanup
  [ -f "ssh/id_rsa_k0s" ] && rm -rf .ssh
}

deleteCluster
createCluster

echo "* Create SSH config"
mkdir -p ~/.ssh
mkdir -p ssh
cp id_rsa_k0s ssh/
cat <<EOF >ssh/config
Host *
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  IdentityFile id_rsa_k0s
  User root
Host controller
  Hostname 127.0.0.1
  Port 9022
Host worker
  Hostname 127.0.0.1
  Port 9023
EOF

echo "* Starting apply"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
echo "* Apply OK"
