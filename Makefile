

SHELL := /bin/bash

curdate=$(shell date --iso-8601='minutes')
build_date = -ldflags "-X  github.com/vogtp/rag/pkg/cfg.BuildInfo=$(curdate)"

# GO_CMD=CGO_ENABLED=0 go
GO_CMD=go

# Branch specific config
BRANCH=$(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')

host=its-a-hack.its.unibas.ch
user=vogtp

.PHONY: run
run: generate ng-build
	$(GO_CMD) run . --log.source web start --log.json --log.level warn | jq -R 'fromjson? | .' 

.PHONY: build
build: generate ng-build
	$(GO_CMD) build $(build_date) -tags prod -o ./build/ . 
	mv ./build/rag ./build/ragctl

.PHONY: generate
generate:
	$(GO_CMD) generate ./...

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
	
.PHONY: pprof-run
pprof-run: build
	nice -19 ./build/ragctl v e path hr-path ./ignore_hr_pdf_failed --pprof.file profile_ragctl --log.level error | tee ignore_pdf_import.log

.PHONY: pprof-cpu
pprof-cpu: 
	go tool pprof -http :1234 ./build/ragctl ignore_profile_ragctl_cpu.pprof

.PHONY: pprof-mem
pprof-mem: 
	go tool pprof -http :1234 ./build/ragctl ignore_profile_ragctl_mem.pprof

.PHONY: pprof-rm
pprof-rm:
	rm ignore_profile_ragctl_*.pprof