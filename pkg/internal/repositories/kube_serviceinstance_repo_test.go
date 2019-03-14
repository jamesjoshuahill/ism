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

package repositories_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeServiceInstanceRepo", func() {
	var (
		repo                    *KubeServiceInstanceRepo
		existingServiceInstance *v1alpha1.ServiceInstance
		resource                types.NamespacedName
	)

	BeforeEach(func() {
		resource = types.NamespacedName{Name: "my-serviceinstance-1", Namespace: "default"}

		existingServiceInstance = &v1alpha1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resource.Name,
				Namespace: resource.Namespace,
			},
			Spec: v1alpha1.ServiceInstanceSpec{
				Name:       "my-serviceinstance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			},
		}

		repo = NewKubeServiceInstanceRepo(kubeClient)
	})

	AfterEach(func() {
		kubeClient.Delete(context.Background(), existingServiceInstance)
	})

	Describe("Get", func() {
		When("the serviceInstance exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingServiceInstance)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns serviceInstance", func() {
				serviceInstance, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(serviceInstance).To(Equal(existingServiceInstance))
			})
		})

		When("the serviceInstance doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(resource)
				Expect(err).To(MatchError("serviceinstances.osbapi.ism.io \"my-serviceinstance-1\" not found"))
			})
		})
	})

	Describe("UpdateStatus", func() {
		When("the serviceInstance exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingServiceInstance)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates status", func() {
				newState := v1alpha1.ServiceInstanceStateProvisioned
				Expect(existingServiceInstance.Status.State).NotTo(Equal(newState))

				err := repo.UpdateState(existingServiceInstance, newState)
				Expect(err).NotTo(HaveOccurred())

				updatedServiceInstance, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatedServiceInstance.Status.State).To(Equal(newState))
				Expect(existingServiceInstance.Status.State).To(Equal(newState))
			})
		})

		When("the serviceInstance doesn't exist", func() {
			It("returns an error", func() {
				newState := v1alpha1.ServiceInstanceStateProvisioned
				err := repo.UpdateState(existingServiceInstance, newState)

				Expect(err).To(MatchError("serviceinstances.osbapi.ism.io \"my-serviceinstance-1\" not found"))
			})
		})
	})
})
