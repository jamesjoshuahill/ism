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

var _ = Describe("KubeBrokerRepo", func() {
	var (
		repo           *KubeBrokerRepo
		existingBroker *v1alpha1.Broker
		resource       types.NamespacedName
	)

	BeforeEach(func() {
		resource = types.NamespacedName{Name: "broker-1", Namespace: "default"}

		existingBroker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resource.Name,
				Namespace: resource.Namespace,
			},
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "http://example.org/broker",
				Username: "admin",
				Password: "password",
			},
		}

		repo = NewKubeBrokerRepo(kubeClient)
	})

	AfterEach(func() {
		kubeClient.Delete(context.Background(), existingBroker)
	})

	Describe("Get", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns broker", func() {
				broker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(broker).To(Equal(existingBroker))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(resource)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})

	Describe("UpdateStatus", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingBroker)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates status", func() {
				newStatus := v1alpha1.BrokerStatus{Registered: &v1alpha1.BrokerStateRegistered{}}
				Expect(existingBroker.Status.Registered).To(BeNil())

				err := repo.UpdateStatus(existingBroker, newStatus)
				Expect(err).NotTo(HaveOccurred())

				updatedBroker, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatedBroker.Status).To(Equal(newStatus))
				Expect(existingBroker.Status).To(Equal(newStatus))
			})
		})

		When("the broker doesn't exist", func() {
			It("returns an error", func() {
				newStatus := v1alpha1.BrokerStatus{Registered: &v1alpha1.BrokerStateRegistered{}}
				err := repo.UpdateStatus(existingBroker, newStatus)

				Expect(err).To(MatchError("brokers.osbapi.ism.io \"broker-1\" not found"))
			})
		})
	})
})
