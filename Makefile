

SHELL := /bin/bash

curdate=$(shell date --iso-8601='minutes')
build_date = -ldflags "-X  github.com/vogtp/rag/pkg/cfg.BuildInfo=$(curdate)"

GO_CMD=CGO_ENABLED=0 go

# Branch specific config
BRANCH=$(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')

host=its-a-hack.its.unibas.ch
user=vogtp

.PHONY: run
run: ng-build
	go run . --log.source web start --log.json | jq

.PHONY: build
build: ng-build
	$(GO_CMD) build $(build_date) -tags prod -o ./build/ . 
	mv ./build/rag ./build/ragctl

.PHONY: ng-build
ng-build:
	cd pkg/web/ng/intrasearch/dist/intrasearch/browser/ ; ng build --base-href=/ui/

.PHONY: remote-stop
remote-stop: remote-stop-rag

.PHONY: remote-stop-%
remote-stop-%:
	ssh root@$(host) systemctl stop $*

.PHONY: remote-start
remote-start: remote-start-rag

.PHONY: remote-start-%
remote-start-%:
	ssh root@$(host) systemctl start $*

.PHONY: remote-restart
remote-restart:	remote-stop-rag remote-start-rag


.PHONY: remote-copy
remote-copy: 
	scp ./build/ragctl $(user)@$(host):srv/rag/

.PHONY: deploy-config
deploy-config: 
	scp ignore_ragctl_intranet.yml $(user)@$(host):srv/rag/

.PHONY: deploy
deploy: build remote-stop-rag remote-copy remote-start-rag remote-autocomplete

.PHONY: remote-autocomplete
remote-autocomplete:
	ssh $(user)@$(host) "srv/rag/ragctl completion bash > ~/.rag.autocomplete"
	ssh $(user)@$(host) "chmod +x ~/.rag.autocomplete"
	