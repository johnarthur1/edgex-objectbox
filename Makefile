#
# Copyright (c) 2018 Cavium
#
# SPDX-License-Identifier: Apache-2.0
#


.PHONY: build clean test docker run

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=CGO_ENABLED=1 GO111MODULE=on go

DOCKERS=docker_build_base docker_volume docker_config_seed docker_export_client docker_export_distro docker_core_data docker_core_metadata docker_core_command docker_support_logging docker_support_notifications docker_sys_mgmt_agent docker_support_scheduler docker_security_secrets_setup docker_security_proxy_setup docker_security_secretstore_setup
.PHONY: $(DOCKERS)

MICROSERVICES=cmd/config-seed/config-seed cmd/export-client/export-client cmd/export-distro/export-distro cmd/core-metadata/core-metadata cmd/core-data/core-data cmd/core-command/core-command cmd/support-logging/support-logging cmd/support-notifications/support-notifications cmd/sys-mgmt-executor/sys-mgmt-executor cmd/sys-mgmt-agent/sys-mgmt-agent cmd/support-scheduler/support-scheduler cmd/security-secrets-setup/security-secrets-setup cmd/security-proxy-setup/security-proxy-setup cmd/security-secretstore-setup/security-secretstore-setup

.PHONY: $(MICROSERVICES)

VERSION=$(shell cat ./VERSION)
DOCKER_TAG=$(VERSION)-dev

#GOFLAGS=-ldflags "-X github.com/objectbox/edgex-objectbox.Version=$(VERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

ARCH=$(shell uname -m)

build: $(MICROSERVICES)

cmd/config-seed/config-seed:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/config-seed

cmd/core-metadata/core-metadata:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/core-metadata

cmd/core-data/core-data:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/core-data

cmd/core-command/core-command:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/core-command

cmd/export-client/export-client:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/export-client

cmd/export-distro/export-distro:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/export-distro

cmd/support-logging/support-logging:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/support-logging

cmd/support-notifications/support-notifications:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/support-notifications

cmd/sys-mgmt-executor/sys-mgmt-executor:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/sys-mgmt-executor

cmd/sys-mgmt-agent/sys-mgmt-agent:
	$(GO) build $(GOFLAGS) -o $@ ./cmd/sys-mgmt-agent

cmd/support-scheduler/support-scheduler:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd/support-scheduler

cmd/security-secrets-setup/security-secrets-setup:
	$(GO) build $(GOFLAGS) -o ./cmd/security-secrets-setup/security-secrets-setup ./cmd/security-secrets-setup

cmd/security-proxy-setup/security-proxy-setup:
	$(GO) build $(GOFLAGS) -o ./cmd/security-proxy-setup/security-proxy-setup ./cmd/security-proxy-setup

cmd/security-secretstore-setup/security-secretstore-setup:
	$(GO) build $(GOFLAGS) -o ./cmd/security-secretstore-setup/security-secretstore-setup ./cmd/security-secretstore-setup


clean:
	rm -f $(MICROSERVICES)

test:
	GO111MODULE=on go test -coverprofile=coverage.out ./...
	GO111MODULE=on go vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-go-mod-tidy.sh
	./bin/test-attribution-txt.sh

run:
	cd bin && ./edgex-launch.sh

run_docker:
	bin/edgex-docker-launch.sh $(EDGEX_DB)

docker: $(DOCKERS)

docker_build_base:
	docker build \
		-f build/base/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-build-base:$(GIT_SHA) \
		.

docker_volume:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-volume:$(GIT_SHA) \
		-t objectboxio/edge-volume:$(DOCKER_TAG) \
		build/volume

docker_config_seed:
	docker build \
		-f cmd/config-seed/Dockerfile \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-core-config-seed:$(GIT_SHA) \
		-t objectboxio/edge-core-config-seed:$(DOCKER_TAG) \
		.

docker_core_metadata:
	docker build \
		-f cmd/Dockerfile --build-arg service=core-metadata \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-core-metadata:$(GIT_SHA) \
		-t objectboxio/edge-core-metadata:$(DOCKER_TAG) \
		.

docker_core_data:
	docker build \
		-f cmd/Dockerfile --build-arg service=core-data \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-core-data:$(GIT_SHA) \
		-t objectboxio/edge-core-data:$(DOCKER_TAG) \
		.

docker_core_command:
	docker build \
		-f cmd/Dockerfile --build-arg service=core-command \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-core-command:$(GIT_SHA) \
		-t objectboxio/edge-core-command:$(DOCKER_TAG) \
		.

docker_export_client:
	docker build \
		-f cmd/Dockerfile --build-arg service=export-client \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-export-client:$(GIT_SHA) \
		-t objectboxio/edge-export-client:$(DOCKER_TAG) \
		.

docker_export_distro:
	docker build \
		-f cmd/export-distro/Dockerfile \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-export-distro:$(GIT_SHA) \
		-t objectboxio/edge-export-distro:$(DOCKER_TAG) \
		.

docker_support_logging:
	docker build \
		-f cmd/support-logging/Dockerfile \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-support-logging:$(GIT_SHA) \
		-t objectboxio/edge-support-logging:$(DOCKER_TAG) \
		.

docker_support_notifications:
	docker build \
		-f cmd/Dockerfile --build-arg service=support-notifications \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-support-notifications:$(GIT_SHA) \
		-t objectboxio/edge-support-notifications:$(DOCKER_TAG) \
		.

docker_support_scheduler:
	docker build \
		-f cmd/Dockerfile --build-arg service=support-scheduler \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-support-scheduler:$(GIT_SHA) \
		-t objectboxio/edge-support-scheduler:$(DOCKER_TAG) \
		.

docker_sys_mgmt_agent:
	docker build \
		-f cmd/sys-mgmt-agent/Dockerfile \
		--build-arg git_sha=$(GIT_SHA) \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-sys-mgmt-agent:$(GIT_SHA) \
		-t objectboxio/edge-sys-mgmt-agent:$(DOCKER_TAG) \
		.

docker_security_secrets_setup:
	# TODO: split this up and rename it when security-secrets-setup is a
	# different container
	docker build \
		-f cmd/security-secrets-setup/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-secret-store:$(GIT_SHA) \
		-t objectboxio/edge-secret-store:$(DOCKER_TAG) \
		.

docker_security_proxy_setup:
	docker build \
		-f cmd/security-proxy-setup/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-security-proxy-setup:$(GIT_SHA) \
		-t objectboxio/edge-security-proxy-setup:$(DOCKER_TAG) \
		.

docker_security_secretstore_setup:
		docker build \
		-f cmd/security-secretstore-setup/Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t objectboxio/edge-security-secretstore-setup:$(GIT_SHA) \
		-t objectboxio/edge-security-secretstore-setup:$(DOCKER_TAG) \
		.

raml_verify:
	docker run --rm --privileged \
		-v $(PWD):/raml-verification -w /raml-verification \
		nexus3.edgexfoundry.org:10003/edgex-docs-builder:$(ARCH) \
		/scripts/raml-verify.sh