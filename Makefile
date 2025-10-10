

SHELL := /bin/bash

curdate=$(shell date --iso-8601='minutes')
build_date = -ldflags "-X  github.com/vogtp/rag/pkg/cfg.BuildInfo=$(curdate)"

# GO_CMD=CGO_ENABLED=0 go
GO_CMD=go

# Branch specific config
BRANCH=$(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')
# prod -> main branch
host_main=its-a-hack.its.unibas.ch
user_main=vogtp
path_main=/srv/intrasearch/
service_main=intrasearch
# qa
host_qm=its-a-hack.its.unibas.ch
user_qm=vogtp
path_qm=/srv/rag/
service_qm=rag
# dev is localhost

user=$(user_$(BRANCH))
path=$(path_$(BRANCH))
host=$(user_$(BRANCH))@$(host_$(BRANCH))
service=$(service_$(BRANCH))

.PHONY: run
run: generate ng-build
	$(GO_CMD) run . --log.source web start --log.json --log.level warn | jq -R 'fromjson? | .' 

.PHONY: build
build: generate ng-build
	$(GO_CMD) build $(build_date) -tags prod -o ./build/ragctl . 

.PHONY: generate
generate:
	$(GO_CMD) generate ./...

.PHONY: ng-serve
ng-serve:
	cd pkg/web/ng/intrasearch/dist/intrasearch/browser/ ; ng serve

.PHONY: ng-build
ng-build:
	cd pkg/web/ng/intrasearch/dist/intrasearch/browser/ ; ng build --base-href=/ui/

.PHONY: remote-stop
remote-stop: remote-stop-$(service)

.PHONY: remote-stop-%
remote-stop-%:
	ssh root@$(host_$(BRANCH)) systemctl stop $*

.PHONY: remote-start
remote-start: remote-start-$(service)

.PHONY: remote-start-%
remote-start-%:
	ssh root@$(host_$(BRANCH)) systemctl start $*

.PHONY: remote-restart
remote-restart:	remote-stop-$(service) remote-start-$(service)


.PHONY: remote-copy
remote-copy: 
	scp ./build/ragctl $(host):$(path_)

.PHONY: remote-copy-config
remote-copy-config: 
	scp ignore_ragctl_$(BRANCH).yml $(host):$(path)/ragctl.yml 

.PHONY: diff-config
diff-config: remote-copy-config
	scp ignore_ragctl_$(BRANCH).yml $(host):$(path)
	ssh $(host) diff $(path)/ignore_ragctl_$(BRANCH).yml $(path)/ragctl.yml 

.PHONY: deploy-config
deploy-config: remote-copy-config remote-restart

.PHONY: deploy
deploy: build remote-stop-rag remote-copy remote-start-rag remote-autocomplete

.PHONY: remote-autocomplete
remote-autocomplete:
	ssh $(host) "$(path_$(BRANCH))ragctl completion bash > ~/.rag.autocomplete"
	ssh $(host) "chmod +x ~/.rag.autocomplete"
	
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