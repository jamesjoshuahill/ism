# ISM

An Independent Services Marketplace built on top of Kubernetes.

### About

There are two parts to ISM - a CLI named `ism` and a corresponding set of Kubernetes CRDs and custom controllers.

### Installation

A [Go](https://golang.org/) development environment is needed in order to build the `ism` CLI.

```
go get github.com/pivotal-cf/ism
cd $GOPATH/github.com/pivotal-cf/ism
make cli
```

Access to a Kubernetes cluster is needed in order to install and run the custom controllers.

```
kubectl cluster-info # ensure you have access to a cluster
make install         # install CRDs into the cluster
make run             # run the custom controllers on your local machine
```

### Usage

```
# display help
./bin/ism --help

# register a broker with ism
./bin/ism broker register \
  --name example-broker \
  --url http://example.broker.com \
  --username x \
  --password y

# list available services and plans
./bin/ism service list

# create an instance of a service
./bin/ism instance create \
  --name example-instance \
  --broker example-broker \
  --service mysql \
  --plan small

# create a binding
./bin/ism binding create \
  --name example-binding \
  --instance example-instance
```

### Development

The following dependencies need to be installed in order to hack on ism:

* [Go](https://golang.org/doc/install)
  * [dep](https://github.com/golang/dep)
  * [ginkgo](https://github.com/onsi/ginkgo)
  * [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)
* [kustomize](https://github.com/kubernetes-sigs/kustomize)
* Access to a k8s cluster

The tests can be run via `make test`.

### Contributing

See [CONTRIBUTING](/CONTRIBUTING.md).
