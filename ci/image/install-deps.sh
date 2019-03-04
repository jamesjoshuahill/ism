#!/usr/bin/env sh

set -e

arch="amd64"

go_version="1.11.5"
dep_version="0.5.0"
kubebuilder_version="1.0.7"
kustomize_version="1.0.11"

# install packages
apt-get -y update
apt-get install -y build-essential git wget

# install golang
wget "https://dl.google.com/go/go${go_version}.linux-${arch}.tar.gz" -O - | tar -C /usr/local -xzf -

# install helpful golang tools
wget -O /usr/local/bin/dep "https://github.com/golang/dep/releases/download/v${dep_version}/dep-linux-${arch}"
chmod +x /usr/local/bin/dep
go get -u github.com/onsi/ginkgo/ginkgo

# install kubebuilder
wget "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${kubebuilder_version}/kubebuilder_${kubebuilder_version}_linux_${arch}.tar.gz" -O - | tar -C /usr/local -xzf -
ln -s "/usr/local/kubebuilder_${kubebuilder_version}_linux_${arch}" "/usr/local/kubebuilder"

# install kustomize
wget "https://github.com/kubernetes-sigs/kustomize/releases/download/v${kustomize_version}/kustomize_${kustomize_version}_linux_amd64" -O "/usr/local/bin/kustomize"
chmod u+x "/usr/local/bin/kustomize"
