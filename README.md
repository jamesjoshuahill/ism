# ISM

An Independent Services Marketplace built on top of Kubernetes.

### About

There are two parts to ISM - a CLI named `ism` and a corresponding set of Kubernetes CRDs and custom controllers.

### Installation

A [Go](https://golang.org/) development environment is needed in order to build the `ism` CLI.

```
go get github.com/pivotal-cf/ism
cd "$(go env GOPATH)/src/github.com/pivotal-cf/ism"
make cli
```

Access to a Kubernetes cluster is needed in order to install and run the custom controllers.

```
kubectl cluster-info # ensure you have access to a cluster
make install         # install CRDs into the cluster
make deploy          # run the custom controllers on the cluster
```

### Usage

```
# display help
ism --help

# register a broker with ism
ism broker register \
  --name example-broker \
  --url http://example.broker.com \
  --username x \
  --password y

# list available services and plans
ism service list

# create an instance of a service
ism instance create \
  --name example-instance \
  --broker example-broker \
  --service mysql \
  --plan small

# create a binding
ism binding create \
  --name example-binding \
  --instance example-instance
  
  
# get a binding (including credentials)
ism binding get \
  --name example-binding

# delete a binding
ism binding delete \
  --name example-binding
```

### Development

The following dependencies need to be installed in order to hack on ism:

* [Go](https://golang.org/doc/install)
  * [dep](https://github.com/golang/dep)
  * [ginkgo](https://github.com/onsi/ginkgo)
  * [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter)
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) (v1.0.8)
* [kustomize](https://github.com/kubernetes-sigs/kustomize)
* Access to a k8s cluster

The tests can be run via `make test`.

You can run the controller locally.
```
make run          # run the custom controllers locally 
```

### Contributing

See [CONTRIBUTING](/CONTRIBUTING.md).
