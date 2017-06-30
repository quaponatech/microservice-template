# Custom settings
VERSION := $(shell git rev-parse --short HEAD)
WHOAMI := $(shell whoami)
OS := $(shell uname -s)
NCPUS := $(shell grep -c ^processor /proc/cpuinfo)

# Dependency management
GITLAB_TOKEN :=

# Certain artifact dependencies
GITLAB_DEPENDENCIES := \
	# "https://gitlab.com/<namespace>/<project>/builds/artifacts/<ref>/download?job=<jobname>"

# The target directories are the copied ones after the artifact has been
# downloaded, see the `deps` build target
# TARGET_DIRS := \
	"include" \
	"lib"

# Go related settings
GO := go
GODOC := godoc
GOMETALINTER := gometalinter
GO_VERBOSE := -v
MOCKGEN := mockgen

# Internally needed paths
SRC_PATH := src
BUILD_PATH := $(subst $(GOPATH)/$(SRC_PATH)/,,$(CURDIR))
LIB_PATH := $(BUILD_PATH)/$(SRC_PATH)/...
CI_PATH := ci
QA_PATH := qa
PROJECT_NAME := $(shell basename $(CURDIR))
UT_COVERAGE_PATH := $(CURDIR)/$(QA_PATH)/ut_coverage
MT_COVERAGE_PATH := $(CURDIR)/$(QA_PATH)/mt_coverage
IT_COVERAGE_PATH := $(CURDIR)/$(QA_PATH)/it_coverage
MOCK_PATH := $(SRC_PATH)/mock_client

# Docker related settings
DOCKER := sudo docker
DEPLOY_DIR := deploy
# TODO: Change to your own registry
REGISTRY := localhost:5000
# TODO: Change this to the specific namespace e.g. 'quapona'
REGISTRY_PATH := $(REGISTRY)/open-source
MAIN := $(DEPLOY_DIR)/main
IMAGE_NAME := $(PROJECT_NAME):$(VERSION)
TARGET_IMAGE_TEST := $(REGISTRY_PATH)/$(IMAGE_NAME)
TARGET_IMAGE_LATEST := $(REGISTRY_PATH)/$(PROJECT_NAME):$(shell git describe --exact-match --tags 2> /dev/null || echo "latest")
IMAGE_PATH := $(DEPLOY_DIR)/$(PROJECT_NAME).tar

# Protobuf related settings
PROTOC := protoc
PROTOBUF_PATH := protobuf

# Rkt related settings
ACBUILD := acbuild

# Kubernetes related settings
# TODO: Change to your kubernetes instance
K8S_SERVER := $(shell minikube ip)
KUBECTL := kubectl
DEPLOYMENT_NAME := $(subst .,-,$(PROJECT_NAME)-$(VERSION))
# TODO: Change to your own chosen service port
SERVICE_PORT := 42302

# Other external tools
GREP := grep

ifneq "$(VERBOSE)" "1"
GO_VERBOSE=
.SILENT:
endif

# Go build targets
all: deps protoc get
	$(GO) build -ldflags "-w -s -linkmode external -extldflags -static"
	file $(PROJECT_NAME) | $(GREP) -q statically
	mv $(PROJECT_NAME) $(MAIN)

assertexecutable:
	if [ ! -e $(MAIN) ]; then \
		echo "Executable '$(MAIN)' does not exist, run 'make all <TARGET>'." ;\
		exit 1 ;\
	fi

get:
	$(GO) get -insecure $(GO_VERBOSE) $(LIB_PATH)

# TODO: Can depend on other projects (when using pb.go as mocked dependency)
mockgen:
	mkdir -p $(MOCK_PATH)
	$(MOCKGEN) $(BUILD_PATH)/$(SRC_PATH)/$(PROTOBUF_PATH) MicroServiceClient > $(MOCK_PATH)/mock_client.go

install: all
	$(GO) install $(GO_VERBOSE) $(LIB_PATH)

update:
	$(GO) get -u $(GO_VERBOSE) $(LIB_PATH)

clean:
	$(GO) clean
	git clean -fdx

# Run source and package documentation server and open it in the browser
doc: all
	-pkill $(GODOC)
	$(GODOC) -http=:6060 -links=true -index -play &
	xdg-open http://localhost:6060/pkg/$(BUILD_PATH)

# Open the coverage results from unit, module and integration tests
# in browser if they exist
coverage:
	$(GO) tool cover -html=$(UT_COVERAGE_PATH).out || true
	$(GO) tool cover -html=$(MT_COVERAGE_PATH).out || true
	$(GO) tool cover -html=$(IT_COVERAGE_PATH).out || true

# Run the benchmarks
bench:
	cd $(SRC_PATH) ;\
	$(GO) test \
		-short \
		$(GO_VERBOSE) \
		-bench=. \
		-cpu=$(NCPUS) \
		-benchtime=1m

# Test help message
test:
	echo "Please chose one of the main test targets:"
	echo "  - utest: Run the unit tests"
	echo "  - mtest: Run the module tests"
	echo "  - itest: Run the integration tests"

# Sets general test environment up.
test_deps:
	mkdir -p $(QA_PATH)

# Run unit tests (see the `-short` flag).
# Prints output to unit test qa file.
# TODO: What does it require? A generated mock of a service?
utest: deps protoc get mockgen test_deps
	cd $(SRC_PATH) ;\
	$(GO) test \
		-short \
		$(GO_VERBOSE) \
		-parallel 8 \
		-timeout 30s \
		-covermode atomic \
		-coverprofile=$(UT_COVERAGE_PATH).out \
		-port=$(SERVICE_PORT) \
		${GOTEST_COLORING}
	$(GO) tool cover -html=$(UT_COVERAGE_PATH).out -o $(UT_COVERAGE_PATH).html
	echo "Unit test coverage results written to $(UT_COVERAGE_PATH).html"

# Run module tests.
# Prints output to module test qa file.
# TODO: This is a kind of test between unit and integration tests.
# 	Extend this to become a real test run if necessary.
# 	Tell also what it requires.
# 	Maybe a generated mock of a service or a running service?
mtest: assertexecutable
	$(MAIN) --dry-run

# TODO: Setup dependencies necessary to run integration tests
integration_deps: dockerload
	# Push image to test registry
	$(DOCKER) tag $(IMAGE_NAME) $(TARGET_IMAGE_TEST)
	$(DOCKER) push $(TARGET_IMAGE_TEST)
	$(DOCKER) rmi -f $(TARGET_IMAGE_TEST)
	# Create a new test environment namespace
	$(KUBECTL) create namespace $(DEPLOYMENT_NAME)
	# Set the current context to that namespace
	$(KUBECTL) config set-context \
		$(shell $(KUBECTL) config current-context) \
		--namespace=$(DEPLOYMENT_NAME)
	# Run the currently built image
	$(KUBECTL) run $(DEPLOYMENT_NAME) \
		--image-pull-policy=Always \
		--image=$(TARGET_IMAGE_TEST) \
		--port=$(SERVICE_PORT)
	$(KUBECTL) expose deploy $(DEPLOYMENT_NAME) \
		--type=NodePort \
		--port=$(SERVICE_PORT) \
		--target-port=$(SERVICE_PORT)
	# TODO: Add/remove needed dependencies
	$(KUBECTL) create -f $(DEPLOY_DIR)/k8s/
	# Give the cluster some time to create the needed PODs
	cnt=1 ; while true; do \
		sleep 3 ;\
		if [ `$(KUBECTL) get pods | tail -n +2 | $(GREP) -v Running | wc -l` != 0 ]; then \
			echo "Waiting for up and running deployment ($$cnt)..." ;\
		else \
			break ;\
		fi ;\
		((cnt = cnt + 1)) ;\
		if [[ $$cnt -ge 11 ]]; then \
			$(KUBECTL) delete namespace $(DEPLOYMENT_NAME) ;\
			exit 1;\
		fi ;\
	done

# Run integration tests (no `-short` flag).
# Prints output to integration test qa file. Requires running database service.
# TODO: What does it require? A running service?
itest: deps protoc get test_deps integration_deps
	cd $(SRC_PATH) ;\
	$(GO) test \
		$(GO_VERBOSE) \
		-parallel 8 \
		-timeout 30s \
		-covermode atomic \
		-coverprofile=$(IT_COVERAGE_PATH).out \
		-ip=$(K8S_SERVER) \
		-port=$(shell $(KUBECTL) get service $(DEPLOYMENT_NAME) | sed -n 's/.*[0-9]\+:\([0-9]\+\)\/TCP.*/\1/p') \
		${GOTEST_COLORING} ;\
		(ret=$$? ;\
		 $(KUBECTL) delete namespace $(DEPLOYMENT_NAME) &&\
		 $(KUBECTL) config set-context $(shell $(KUBECTL) config current-context) --namespace=default &&\
		 exit $$ret)
	$(GO) tool cover -html=$(IT_COVERAGE_PATH).out -o $(IT_COVERAGE_PATH).html
	echo "Integration test coverage results written to $(IT_COVERAGE_PATH).html"

lint: deps protoc get
	$(GOMETALINTER) \
		--enable-all \
		--line-length=120 \
		--deadline=60s \
		--exclude=$(MOCK_PATH) \
		--exclude=$(SRC_PATH)/$(PROTOBUF_PATH) \
		./...

# Protocol buffers build targets
protoc:
	$(PROTOC) --go_out=plugins=grpc:$(SRC_PATH) $(PROTOBUF_PATH)/*.proto

# Build docker and ACI images
images: clean docker rkt

assertdockerimage:
	if [ ! -e $(IMAGE_PATH) ]; then \
		echo "Docker image '$(IMAGE_PATH)' does not exist, run 'make all docker <TARGET>'." ;\
		exit 1 ;\
	fi

# Create a docker image as a local file
docker: assertexecutable
	$(DOCKER) build --pull -t $(IMAGE_NAME) $(DEPLOY_DIR)
	$(DOCKER) save $(IMAGE_NAME) -o $(IMAGE_PATH)
	sudo chown -R $(WHOAMI) $(DEPLOY_DIR)
	$(DOCKER) rmi -f $(IMAGE_NAME)

# Loads the local docker image into the registry
dockerload: assertdockerimage
	$(DOCKER) load -i $(IMAGE_PATH)

# Publish the local docker image as 'latest'
dockerpublish: docker dockerload
	$(DOCKER) tag $(IMAGE_NAME) $(TARGET_IMAGE_LATEST)
	$(DOCKER) push $(TARGET_IMAGE_LATEST)
	$(DOCKER) rmi -f $(IMAGE_NAME) $(TARGET_IMAGE_LATEST)

dockerclean:
	$(DOCKER) images -f=dangling=true -qa | xargs -r $(DOCKER) rmi -f

# Rkt build targets
rkt: assertexecutable
	cd $(DEPLOY_DIR) && $(ACBUILD) script Rkt.acb

# Dependency download from external URLs
deps:
	set -e ;\
	c=1 ;\
	for directory in $(TARGET_DIRS); do \
		mkdir -p vendor/$$directory ;\
	done ;\
	for dep in $(DEPENDENCIES); do \
		target=vendor/tmp/$$c ;\
		if [ ! -d $$target ]; then \
			mkdir -p $$target ;\
			quiet=-qq ;\
			silent=-s ;\
			if [ ! -z "$(VERBOSE)" ]; then \
				echo "Downloading '$$dep'..." ;\
				quiet= ;\
				silent= ;\
			fi ;\
			curl -skLH PRIVATE-TOKEN:$(TOKEN) $$dep -o $$target/$$c.zip ;\
			if [ ! -z "$(VERBOSE)" ]; then \
				echo "Extracting" ;\
			fi ;\
			pushd $$target > /dev/null ;\
			unzip $$quiet $$c.zip ;\
			for directory in $(TARGET_DIRS); do \
				find . -type d -name $$directory -exec cp -r {}/. ../../$$directory \; ;\
			done ;\
			popd > /dev/null ;\
		else \
			if [ ! -z "$(VERBOSE)" ]; then \
				echo "Skipping '$$dep' since '$$target' already exists" ;\
			fi ;\
		fi ;\
		((c=c+1));\
	done
