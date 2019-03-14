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
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("KubePlanRepo", func() {
	var (
		repo              *KubePlanRepo
		brokerService     *v1alpha1.BrokerService
		brokerServicePlan *v1alpha1.BrokerServicePlan
	)

	BeforeEach(func() {
		brokerService = &v1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "service-id-1",
				Namespace: "default",
			},
		}

		brokerServicePlan = &v1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "plan-id-1",
				Namespace: "default",
			},
		}

		repo = NewKubePlanRepo(kubeClient)
	})

	Describe("Create", func() {
		When("the service exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), brokerService)
				Expect(err).NotTo(HaveOccurred())

				err = repo.Create(brokerService, osbapi.Plan{
					ID:          "plan-id-1",
					Name:        "plan-1-name",
					Description: "plan-1-description",
				})
				Expect(err).NotTo(HaveOccurred())

				err = kubeClient.Get(context.Background(), types.NamespacedName{Name: "plan-id-1", Namespace: "default"}, brokerServicePlan)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(kubeClient.Delete(context.Background(), brokerService)).To(Succeed())
				Expect(kubeClient.Delete(context.Background(), brokerServicePlan)).To(Succeed())
			})

			It("creates the plan with the correct spec", func() {
				Expect(brokerServicePlan.Spec).To(Equal(v1alpha1.BrokerServicePlanSpec{
					Name: "plan-1-name",
				}))
			})

			It("generates the correct name and namespace", func() {
				Expect(brokerServicePlan.ObjectMeta.Name).To(Equal("plan-id-1"))
				Expect(brokerServicePlan.ObjectMeta.Namespace).To(Equal("default"))
			})

			It("sets the owner reference of the plan to the service", func() {
				Expect(brokerServicePlan.ObjectMeta.OwnerReferences).To(HaveLen(1))
				Expect(brokerServicePlan.ObjectMeta.OwnerReferences[0].UID).To(Equal(brokerService.ObjectMeta.UID))
			})

			When("the service is invalid", func() {
				It("returns an error", func() {
					invalidService := &v1alpha1.BrokerService{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "service-without-uid",
						},
					}

					err := repo.Create(invalidService, osbapi.Plan{
						ID:          "plan-id-1",
						Name:        "plan-1-name",
						Description: "plan-1-description",
					})

					Expect(err).To(MatchError("BrokerServicePlan.osbapi.ism.io \"plan-id-1\" is invalid: " +
						"metadata.ownerReferences.uid: Invalid value: \"\": uid must not be empty"))
				})
			})
		})
	})
})
