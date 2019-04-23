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
	"time"

	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/repositories/kube"
)

var _ = Describe("Plan", func() {
	var (
		kubeClient client.Client
		planRepo   *Plan
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		planRepo = &Plan{
			KubeClient: kubeClient,
		}
	})

	Describe("GetPlanByID", func() {
		var (
			plan *osbapi.Plan
			err  error
		)

		JustBeforeEach(func() {
			plan, err = planRepo.GetPlanByID("plan-1")
		})

		When("the plan exists", func() {
			BeforeEach(func() {
				planResource := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name:      "my-plan",
						ServiceID: "service-1",
					},
				}

				Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())
			})

			AfterEach(func() {
				deletePlans(kubeClient, "plan-1")
			})

			It("returns the plan", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(plan).To(Equal(&osbapi.Plan{
					ID:        "plan-1",
					Name:      "my-plan",
					ServiceID: "service-1",
				}))
			})
		})

		When("the plan does not exist", func() {
			It("returns a not found error", func() {
				Expect(err).To(MatchError("plan not found"))
			})
		})

		When("the client errors", func() {
			BeforeEach(func() {
				// The purpose of this test is to ensure that we are properly propagating errors from the kubeclient

				// We change the config to be unreachable after the client has been created because
				// the client checks the connection on client.New
				unreachableCfg := *cfg
				badKubeClient, clientErr := client.New(&unreachableCfg, client.Options{Scheme: scheme.Scheme})
				Expect(clientErr).NotTo(HaveOccurred())

				unreachableCfg.Host = "192.0.2.1"
				unreachableCfg.Timeout = time.Second

				planRepo = &Plan{
					KubeClient: badKubeClient,
				}
			})

			It("propagates the error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetPlans", func() {
		var (
			plans []*osbapi.Plan
			err   error
		)

		JustBeforeEach(func() {
			plans, err = planRepo.GetPlans("service-1")
		})

		When("plans contain owner references to services", func() {
			BeforeEach(func() {
				planResource := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-1",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "service-1",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "service-uid-1",
						}},
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())

				planResource2 := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-2",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "service-2",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "service-uid-2",
						}},
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan-2",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource2)).To(Succeed())
			})

			AfterEach(func() {
				deletePlans(kubeClient, "plan-1", "plan-2")
			})

			It("returns plans by service id", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(plans).To(Equal([]*osbapi.Plan{{
					Name:      "my-plan",
					ID:        "plan-1",
					ServiceID: "service-1",
				}}))
			})
		})

		When("the plan owner reference is not set", func() {
			BeforeEach(func() {
				planResource := &v1alpha1.BrokerServicePlan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "plan-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServicePlanSpec{
						Name: "my-plan",
					},
				}
				Expect(kubeClient.Create(context.TODO(), planResource)).To(Succeed())
			})

			AfterEach(func() {
				deletePlans(kubeClient, "plan-1")
			})

			It("successfully returns no plans", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(plans).To(HaveLen(0))
			})
		})
	})
})

func deletePlans(kubeClient client.Client, planNames ...string) {
	for _, p := range planNames {
		pToDelete := &v1alpha1.BrokerServicePlan{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), pToDelete)).To(Succeed())
	}
}
