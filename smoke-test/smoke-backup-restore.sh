#!/usr/bin/env sh

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl.yaml"}

set -e

. ./smoke.common.sh
trap runCleanup EXIT

# custom exit trap to cleanup the backup archives
runCleanup() {
    cleanup
    rm k0s_backup*.tar.gz
}

deleteCluster
createCluster
../cfctl init
../cfctl apply --config "${CFCTL_CONFIG}" --debug

# Collect some facts so we can validate restore actually did full restore
system_ns_uid=$(bootloose ssh root@manager0 -- k0s kubectl --kubeconfig "/var/lib/k0s/pki/admin.conf" get -n kube-system namespace kube-system -o template='{{.metadata.uid}}')
node_uid=$(bootloose ssh root@manager0 -- k0s kubectl --kubeconfig "/var/lib/k0s/pki/admin.conf" get node worker0 -o template='{{.metadata.uid}}')

../cfctl backup --config "${CFCTL_CONFIG}" --debug

# Reset the controller
bootloose ssh root@manager0 -- k0s stop
bootloose ssh root@manager0 -- k0s reset

../cfctl apply --config "${CFCTL_CONFIG}" --debug --restore-from $(ls k0s_backup*.tar.gz)

# Verify kube object UIDs match so we know we did full restore of the API objects
new_system_ns_uid=$(bootloose ssh root@manager0 -- k0s kubectl --kubeconfig "/var/lib/k0s/pki/admin.conf" get -n kube-system namespace kube-system -o template='{{.metadata.uid}}')
if [ "$system_ns_uid" != "$new_system_ns_uid" ]; then
    echo "kube-system UIDs do not match after restore!!!"
    exit 1
fi
new_node_uid=$(bootloose ssh root@manager0 -- k0s kubectl --kubeconfig "/var/lib/k0s/pki/admin.conf" get node worker0 -o template='{{.metadata.uid}}')
if [ "$node_uid" != "$new_node_uid" ]; then
    echo "worker0 UIDs do not match after restore!!!"
    exit 1
fi
