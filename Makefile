
# Image URL to use all building/pushing image targets
IMG ?= ismteam/controller:latest
BROKERIMG ?= mattmcneeney/overview-broker:latest
CLI_NAME = bin/ism
GINKGO_ARGS = -r -p -randomizeSuites -randomizeAllSpecs

all: clean test manager cli

# Run tests
test: fmt vet manifests unit-tests integration-tests acceptance-tests

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager github.com/pivotal-cf/ism/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
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
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate:
	go generate ./...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

### CUSTOM MAKE RULES ###

deploy-test-broker:
	kubectl apply -f acceptance/assets/broker

run-test-broker:
	docker run -d -p 127.0.0.1:1122:8080/tcp mattmcneeney/overview-broker

stop-test-broker:
	docker ps | grep mattmcneeney/overview-broker | awk '{print $$1}' | xargs -n1 docker kill

cli:
	go build -o ${CLI_NAME} cmd/ism/main.go

clean:
	rm -f ${CLI_NAME}

# Cannot yet -randomizeAllSpecs the acceptance tests
acceptance-tests:
	ginkgo -r acceptance

# skip integration/acceptance tests
unit-tests:
	ginkgo ${GINKGO_ARGS} -skipPackage acceptance,kube,pkg/controller,pkg/api,pkg/internal/repositories

integration-tests: cli-integration-tests kube-integration-tests

cli-integration-tests:
	ginkgo ${GINKGO_ARGS} kube

kube-integration-tests:
	ginkgo ${GINKGO_ARGS} pkg/controller pkg/api pkg/internal/repositories

delete-test-broker:
	kubectl delete -f acceptance/assets/broker

delete-controller:
	kustomize build config/default | kubectl delete -f -

uninstall-crds:
	kubectl delete -f config/crds
