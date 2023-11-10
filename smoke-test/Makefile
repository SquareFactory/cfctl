
bootloose := $(shell which bootloose)
ifeq ($(bootloose),)
bootloose := $(shell go env GOPATH)/bin/bootloose
endif

envsubst := $(shell which envsubst)
ifeq ($(envsubst),)
$(error 'envsubst' NOT found in path, please install it and re-run)
endif

.PHONY: cfctl
cfctl:
	$(MAKE) -C .. cfctl

$(bootloose):
	go install github.com/k0sproject/bootloose@latest

id_rsa_k0s:
	ssh-keygen -t rsa -f ./id_rsa_k0s -N ""

smoke-basic: $(bootloose) id_rsa_k0s cfctl
	./smoke-basic.sh

smoke-basic-openssh: $(bootloose) id_rsa_k0s cfctl
	./smoke-basic-openssh.sh

smoke-dynamic: $(bootloose) id_rsa_k0s cfctl
	./smoke-dynamic.sh

smoke-files: $(bootloose) id_rsa_k0s cfctl
	./smoke-files.sh

smoke-init: $(bootloose) id_rsa_k0s cfctl
	./smoke-init.sh

smoke-upgrade: $(bootloose) id_rsa_k0s cfctl
	./smoke-upgrade.sh

smoke-dryrun: $(bootloose) id_rsa_k0s cfctl
	./smoke-dryrun.sh

smoke-reset: $(bootloose) id_rsa_k0s cfctl
	./smoke-reset.sh

smoke-os-override: $(bootloose) id_rsa_k0s cfctl
	BOOTLOOSE_TEMPLATE=bootloose.yaml.osoverride.tpl CFCTL_CONFIG=cfctl-single.yaml OS_RELEASE_PATH=$(realpath os-release) OS_OVERRIDE="ubuntu" ./smoke-basic.sh

smoke-backup-restore: $(bootloose) id_rsa_k0s cfctl
	./smoke-backup-restore.sh

%.iid: Dockerfile.%
	docker build --iidfile '$@' - < '$<'
