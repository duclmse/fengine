.PHONY: all $(SERVICES) dockers dockers_dev latest release

BUILD_DIR = build

SERVICES = \
	fengine
#	devices devicetypes groups httpgw msggw mqttv2 devicesauthn devicetokens topics msgwriter httpmetricscollector \
#	devicevars organizations users eventaction varsparser eventhandler roles events attributes zigbeeadapter pricing \

DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
TEST = $(addprefix test_,$(SERVICES))
RELEASE = $(addprefix release_,$(SERVICES))
PLATFORM = vht-iot-platform
REPOSITORY = 172.20.1.22
CGO_ENABLED ?= 0
GOARCH ?= amd64
GOOS ?= linux

$(SERVICES):
	@echo "> Compiling '$(SERVICES)'..."
	$(call compile_service,$(@))

$(TEST):
	@echo "> Testing '$(TEST)'..."
	$(call test_service,$(@))

$(DOCKERS):
	@echo "> Building docker '$(DOCKERS)'..."
	$(call make_docker,$(@),$(GOARCH))

$(DOCKERS_DEV):
	@echo "> Building dev docker '$(DOCKERS_DEV)'..."
	$(call make_docker_dev,$(@))

$(RELEASE):
	@echo "> Releasing '$(RELEASE)'..."
	$(call make_release_svc,$(@))

all: $(SERVICES)

clean:
	rm -rf ${BUILD_DIR}

clean_docker:
	docker-compose -f docker/docker-compose.yml stop
	docker ps -f name=$(PLATFORM) -aq | xargs -r docker rm
	docker ps -f name=$(PLATFORM) -f status=dead -f status=exited -aq | xargs -r docker rm -v
	docker images "$(PLATFORM)\/*" -f dangling=true -q | xargs -r docker rmi
	docker images -q $(PLATFORM)\/* | xargs -r docker rmi
ifdef pv
	docker volume ls -f name=$(PLATFORM) -f dangling=true -q | xargs -r docker volume rm
endif

install:
	cp $(BUILD_DIR)/* $(GOBIN)

test:
	#go test -mod=vendor -v -race -count 1 -tags test $(shell go list ./... | grep -v 'vendor\|cmd')
	go test -mod=vendor -v -race -count 1 -coverprofile .coverage.txt \
		-tags test $(shell go list ./... | grep -v 'vendor\|cmd')
	go tool cover -func .coverage.txt

# sudo apt-get install gogoprotobuf
proto:
	@protoc --gofast_out=plugins=grpc:. pb/*.proto
	@echo "Done generating"

# sudo npm i -g grpc grpc-tools grpc_tools_node_protoc_ts
jspb:
	@grpc_tools_node_protoc --js_out=import_style=commonjs,binary:executor/src/ --ts_out=executor/src/ \
		--grpc_out=grpc_js:executor/src/ --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts pb/*.proto
	@awk '{sub(/import \* as grpc from "grpc";/, "import * as grpc from \"@grpc/grpc-js\";"); print}' \
		./executor/src/pb/fengine_grpc_pb.d.ts > t && mv t ./executor/src/pb/fengine_grpc_pb.d.ts
	@echo "Done generating"

dockers: $(DOCKERS)

dockers_dev: $(DOCKERS_DEV)

changelog:
	git log $(shell git describe --tags --abbrev=0)..HEAD --pretty=format:"- %s"

latest: dockers
	$(call docker_push,latest)

release:
	$(eval VERSION=$(shell git describe --abbrev=0 --tags | cut -d '_' -f1))
	@echo $(VERSION)
	$(MAKE) dockers
	for svc in $(SERVICES); do \
		docker tag $(REPOSITORY)/$(PLATFORM)/$$svc $(REPOSITORY)/$(PLATFORM)/$$svc:$(VERSION); \
	done
	$(call docker_push,$(VERSION))
	for svc in $(SERVICES); do \
		docker rmi $(REPOSITORY)/$(PLATFORM)/$$svc ; \
		docker rmi $(REPOSITORY)/$(PLATFORM)/$$svc:$(VERSION) ; \
	done

rundev:
	cd scripts && ./run.sh

run:
	docker-compose -f docker/docker-compose.yml up

# Run all core services except distributed tracing system - Jaeger. Recommended on gateways:
rungw:
	MF_JAEGER_URL= docker-compose -f docker/docker-compose.yml up --scale jaeger=0

push:
	${call docker_push}

define docker_push
	@for svc in $(SERVICES); do \
		docker push $(REPOSITORY)/$(PLATFORM)/$$svc; \
	done
endef

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
		go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/viot-$(1) cmd/$(1)/main.go
endef

define test_service
	$(eval svc=$(subst test_,,$(1)))
	go test -mod=vendor -v -race -count 1 -tags test $(shell go list ./$(svc))
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(REPOSITORY)/$(PLATFORM)/$(svc) \
		-f docker/Dockerfile .
	docker image prune -f --filter label=stage=builder
endef

define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))

	docker build --no-cache --build-arg SVC=$(svc) --tag=$(PLATFORM)/$(svc) -f docker/Dockerfile.dev ./build
endef

define make_release_svc
	$(eval svc=$(subst release_,,$(1)))
	$(call make_docker, docker_$(svc))
	docker tag  $(REPOSITORY)/$(PLATFORM)/$(svc) $(REPOSITORY)/$(PLATFORM)/$(svc):$(version)
	docker push $(REPOSITORY)/$(PLATFORM)/$(svc):$(version)
	docker rmi  $(REPOSITORY)/$(PLATFORM)/$(svc)
	docker rmi  $(REPOSITORY)/$(PLATFORM)/$(svc):$(version)
endef
