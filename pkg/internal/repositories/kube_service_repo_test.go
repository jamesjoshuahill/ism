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
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeServiceRepo", func() {
	var (
		repo            *KubeServiceRepo
		broker          *v1alpha1.Broker
		brokerService   *v1alpha1.BrokerService
		returnedService *v1alpha1.BrokerService
	)

	BeforeEach(func() {
		broker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-broker",
				Namespace: "default",
			},
		}

		repo = NewKubeServiceRepo(kubeClient)
	})

	Describe("Create", func() {
		When("the broker exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), broker)
				Expect(err).NotTo(HaveOccurred())

				returnedService, err = repo.Create(broker, osbapi.Service{
					ID:          "service-id-1",
					Name:        "service-one",
					Description: "service-description",
				})
				Expect(err).NotTo(HaveOccurred())

				brokerService = &v1alpha1.BrokerService{}
				err = kubeClient.Get(context.Background(), types.NamespacedName{Name: "service-id-1", Namespace: "default"}, brokerService)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(kubeClient.Delete(context.Background(), brokerService)).To(Succeed())
				Expect(kubeClient.Delete(context.Background(), broker)).To(Succeed())
			})

			It("returns the created service", func() {
				Expect(returnedService.Spec).To(Equal(brokerService.Spec))
			})

			It("creates the service with the correct spec", func() {
				Expect(brokerService.Spec).To(Equal(v1alpha1.BrokerServiceSpec{
					Name:        "service-one",
					Description: "service-description",
					BrokerName:  "my-broker",
				}))
			})

			It("generates the correct name and namespace", func() {
				Expect(brokerService.ObjectMeta.Name).To(Equal("service-id-1"))
				Expect(brokerService.ObjectMeta.Namespace).To(Equal("default"))
			})

			It("sets the owner reference of the service to the broker", func() {
				Expect(brokerService.ObjectMeta.OwnerReferences).To(HaveLen(1))
				Expect(brokerService.ObjectMeta.OwnerReferences[0].UID).To(Equal(broker.ObjectMeta.UID))
			})
		})

		When("the broker is invalid", func() {
			It("returns an error", func() {
				invalidBroker := &v1alpha1.Broker{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Name:      "broker-without-uid",
					},
				}

				_, err := repo.Create(invalidBroker, osbapi.Service{
					ID:          "service-id-1",
					Name:        "service-one",
					Description: "service-description",
				})

				Expect(err).To(MatchError("BrokerService.osbapi.ism.io \"service-id-1\" is invalid" +
					": metadata.ownerReferences.uid: Invalid value: \"\": uid must not be empty"))
			})
		})
	})
})
