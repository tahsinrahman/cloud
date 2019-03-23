
# Image URL to use all building/pushing image targets
IMG ?= controller:latest

all: test manager

# Run tests
test: verify generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/pharmer/cloud/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate: clientset
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

.PHONY: clientset
clientset: ## Generate a typed clientset
	rm -rf pkg/client
	cd ./vendor/k8s.io/code-generator/cmd && go install ./client-gen ./lister-gen ./informer-gen
	$$GOPATH/bin/client-gen --clientset-name clientset --input-base github.com/pharmer/cloud/pkg/apis \
		--input cloud/v1 --output-package github.com/pharmer/cloud/pkg/client/clientset_generated \
		--go-header-file=./hack/boilerplate.go.txt

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

.PHONY: verify
verify:
	./hack/verify_clientset.sh
