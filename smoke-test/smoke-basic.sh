#!/bin/bash

CFCTL_CONFIG=${CFCTL_CONFIG:-"cfctl.yaml"}

set -e

. ./smoke.common.sh
trap cleanup EXIT

deleteCluster
createCluster

echo "* Starting apply"
../cfctl apply --config "${CFCTL_CONFIG}" --debug
echo "* Apply OK"

echo "* Verify hooks were executed on the host"
footloose ssh root@manager0 -- grep -q hello apply.hook

echo "* Verify 'cfctl kubeconfig' output includes 'data' block"
../cfctl kubeconfig --config cfctl.yaml | grep -v -- "-data"

echo "* Run kubectl on controller"
footloose ssh root@manager0 -- k0s kubectl get nodes

echo "* Downloading kubectl for local test"
downloadKubectl

echo "* Using cfctl kubecofig locally"
../cfctl kubeconfig --config cfctl.yaml >kubeconfig

echo "* Output:"
cat kubeconfig | grep -v -- "-data"

echo "* Running kubectl"
./kubectl --kubeconfig kubeconfig get nodes
echo "* Done"
