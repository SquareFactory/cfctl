
footloose := $(shell which footloose)
ifeq ($(footloose),)
footloose := $(shell go env GOPATH)/bin/footloose
endif

envsubst := $(shell which envsubst)
ifeq ($(envsubst),)
$(error 'envsubst' NOT found in path, please install it and re-run)
endif

.PHONY: cfctl
cfctl:
	$(MAKE) -C .. cfctl

$(footloose):
	go install github.com/weaveworks/footloose/...@0.6.3

id_rsa_k0s:
	ssh-keygen -t rsa -f ./id_rsa_k0s -N ""

smoke-basic: $(footloose) id_rsa_k0s cfctl
	./smoke-basic.sh

smoke-dynamic: $(footloose) id_rsa_k0s cfctl
	./smoke-dynamic.sh

smoke-files: $(footloose) id_rsa_k0s cfctl
	./smoke-files.sh

smoke-init: $(footloose) id_rsa_k0s cfctl
	./smoke-init.sh

smoke-upgrade: $(footloose) id_rsa_k0s cfctl
	./smoke-upgrade.sh

smoke-reset: $(footloose) id_rsa_k0s cfctl
	./smoke-reset.sh

smoke-os-override: $(footloose) id_rsa_k0s cfctl
	FOOTLOOSE_TEMPLATE=footloose.yaml.osoverride.tpl CFCTL_CONFIG=cfctl-single.yaml OS_RELEASE_PATH=$(realpath os-release) OS_OVERRIDE="ubuntu" ./smoke-basic.sh

smoke-backup-restore: $(footloose) id_rsa_k0s cfctl
	./smoke-backup-restore.sh
