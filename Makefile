
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

all: pharmer-tools

# Run tests
test: generate manifests fmt vet
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build pharmer-tools binary
pharmer-tools: generate fmt vet
	go build -o bin/pharmer-tools github.com/pharmer/cloud/cmd/pharmer-tools

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/pharmer-tools/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crd
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook object paths="./pkg/apis/..." output:crd:artifacts:config=config/crd

# Run go fmt against code
fmt:
	gofmt -s -w ./pkg ./cmd
	goimports -w ./pkg ./cmd

# Run go vet against code
vet:
	# go vet ./pkg/... ./cmd/...

# Generate code
.PHONY: clientset
clientset: client-gen ## Generate a typed clientset
	rm -rf pkg/client
	$(CLIENT_GEN) --clientset-name clientset --input-base github.com/pharmer/cloud/pkg/apis \
		--input cloud/v1 --output-package github.com/pharmer/cloud/pkg/client/clientset_generated \
		--go-header-file=./hack/boilerplate.go.txt

generate: controller-gen clientset
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/...

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0-alpha.2
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

# find or download controller-gen
# download client-gen if necessary
client-gen:
ifeq (, $(shell which client-gen))
	go get k8s.io/code-generator/cmd/client-gen@639c964206c28ac3859cf36f212c24775616884a
CLIENT_GEN=$(shell go env GOPATH)/bin/client-gen
else
CLIENT_GEN=$(shell which client-gen)
endif
