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
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/repositories/kube"
)

var _ = Describe("Binding", func() {

	var (
		kubeClient client.Client

		bindingRepo *Binding
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		bindingRepo = &Binding{
			KubeClient: kubeClient,
		}
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			b := &osbapi.Binding{
				Name:       "my-binding",
				InstanceID: "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			}

			err = bindingRepo.Create(b)
		})

		AfterEach(func() {
			deleteBindings(kubeClient, "my-binding")
		})

		It("creates a new ServiceBinding resource", func() {
			Expect(err).NotTo(HaveOccurred())

			key := types.NamespacedName{
				Name:      "my-binding",
				Namespace: "default",
			}

			fetched := &v1alpha1.ServiceBinding{}
			Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

			Expect(fetched.Spec).To(Equal(v1alpha1.ServiceBindingSpec{
				Name:       "my-binding",
				InstanceID: "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			}))
		})

		When("creating a new Binding fails", func() {
			BeforeEach(func() {
				// create the binding first, so that the second create errors
				b := &osbapi.Binding{
					ID:         "binding-1",
					Name:       "my-binding",
					InstanceID: "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "broker-1",
				}

				Expect(bindingRepo.Create(b)).To(Succeed())
			})

			It("propagates the error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

func deleteBindings(kubeClient client.Client, bindingNames ...string) {
	for _, b := range bindingNames {
		bToDelete := &v1alpha1.ServiceBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      b,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
	}
}
