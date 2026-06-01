GO_VERSION     := 1.26
ALPINE_VERSION := 3.19.0

VERSION         := 0.0.1
DEVELOP_VERSION := develop-${VERSION}
GIT_COMMIT      := $(shell git log -1 --pretty=format:%h)

STORE_IMAGE           := store:${DEVELOP_VERSION}
SENDER_IMAGE          := sender:${DEVELOP_VERSION}
QUERIER_IMAGE         := querier:${DEVELOP_VERSION}
TENANT_REGISTRY_IMAGE := tenant-registry:${DEVELOP_VERSION}

develop: build compose

build: build-store build-sender build-querier build-tenant-registry

build-store:
	docker build -f store.Dockerfile \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
    	--build-arg VERSION=${DEVELOP_VERSION} \
    	--build-arg GIT_COMMIT=${GIT_COMMIT} \
    	--tag ${STORE_IMAGE} .

build-sender:
	docker build -f sender.Dockerfile \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
		--build-arg VERSION=${DEVELOP_VERSION} \
		--build-arg GIT_COMMIT=${GIT_COMMIT} \
		--tag ${SENDER_IMAGE} .

build-querier:
	docker build -f querier.Dockerfile \
    	--build-arg GO_VERSION=${GO_VERSION} \
    	--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
    	--build-arg VERSION=${DEVELOP_VERSION} \
    	--build-arg GIT_COMMIT=${GIT_COMMIT} \
    	--tag ${QUERIER_IMAGE} .

build-tenant-registry:
	docker build -f tenant-registry.Dockerfile \
    	--build-arg GO_VERSION=${GO_VERSION} \
    	--build-arg ALPINE_VERSION=${ALPINE_VERSION} \
    	--build-arg VERSION=${DEVELOP_VERSION} \
    	--build-arg GIT_COMMIT=${GIT_COMMIT} \
    	--tag ${TENANT_REGISTRY_IMAGE} .

compose:
	STORE_IMAGE=${STORE_IMAGE} \
	SENDER_IMAGE=${SENDER_IMAGE} \
	QUERIER_IMAGE=${QUERIER_IMAGE} \
	TENANT_REGISTRY_IMAGE=${TENANT_REGISTRY_IMAGE} \
	docker compose -f develop/docker-compose.yml up --force-recreate --remove-orphans

