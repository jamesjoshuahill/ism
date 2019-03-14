/*
Copyright (C) 2019-Present Pivotal Software, Inc. All rights reserved.

This program and the accompanying materials are made available under the terms
of the under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package kube_test

import (
	"path/filepath"
	"testing"

	"github.com/pivotal-cf/ism/pkg/apis"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var cfg *rest.Config

func TestKube(t *testing.T) {
	var testEnv *envtest.Environment

	BeforeSuite(func() {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "config", "crds")},
		}
		apis.AddToScheme(scheme.Scheme)

		var err error
		cfg, err = testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		testEnv.Stop()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Kube Suite")
}
