SHELL=/bin/bash -o pipefail

GO_PKG   := pharmer.dev
REPO     := $(notdir $(shell pwd))

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

BUILD_IMAGE ?= appscode/gengo:release-1.14

all: pharmer-tools

# Run tests
test: manifests fmt vet
	go test -mod=vendor ./pkg/... ./cmd/... -coverprofile cover.out

# Build pharmer-tools binary
pharmer-tools: fmt vet
	go build -mod=vendor -o bin/pharmer-tools $(GO_PKG)/$(REPO)/cmd/pharmer-tools

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run -mod=vendor ./cmd/pharmer-tools/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd

# Run go fmt against code
fmt:
	gofmt -s -w ./pkg ./cmd
	goimports -w ./pkg ./cmd

# Run go vet against code
vet:
	# go vet ./pkg/... ./cmd/...

DOCKER_REPO_ROOT := /go/src/$(GO_PKG)/$(REPO)

# Generate a typed clientset
.PHONY: clientset
clientset:
	@rm -rf pkg/client
	@docker run --rm -ti                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
	    --env HTTP_PROXY=$(HTTP_PROXY)                   \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(BUILD_IMAGE)                                   \
		/go/src/k8s.io/code-generator/generate-groups.sh \
			"deepcopy,client"                            \
			$(GO_PKG)/$(REPO)/pkg/client                 \
			$(GO_PKG)/$(REPO)/pkg/apis                   \
			"cloud:v1"                                   \
			--go-header-file "./hack/boilerplate.go.txt"

# Generate CRD manifests
.PHONY: manifests
manifests:
	@echo "Generating CRD manifests"
	@docker run --rm -ti                    \
		-u $$(id -u):$$(id -g)              \
		-v /tmp:/.cache                     \
		-v $$(pwd):$(DOCKER_REPO_ROOT)      \
		-w $(DOCKER_REPO_ROOT)              \
	    --env HTTP_PROXY=$(HTTP_PROXY)      \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)    \
		$(BUILD_IMAGE)                      \
		controller-gen                      \
			$(CRD_OPTIONS)                  \
			paths="./pkg/apis/..."          \
			output:crd:artifacts:config=config/crd

.PHONY: gen
gen: clientset manifests
