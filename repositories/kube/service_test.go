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
	. "github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/repositories/kube"
)

var _ = Describe("Service", func() {
	var (
		kubeClient  client.Client
		serviceRepo *Service
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		serviceRepo = &Service{
			KubeClient: kubeClient,
		}
	})

	Describe("GetServiceByID", func() {
		var (
			service *osbapi.Service
			err     error
		)

		JustBeforeEach(func() {
			service, err = serviceRepo.GetServiceByID("service-1")
		})

		When("the service exists", func() {
			BeforeEach(func() {
				serviceResource := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service",
						Description: "my-service-description",
						BrokerName:  "my-broker",
					},
				}

				Expect(kubeClient.Create(context.TODO(), serviceResource)).To(Succeed())
			})

			AfterEach(func() {
				deleteServices(kubeClient, "service-1")
			})

			It("returns the service", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(service).To(Equal(&osbapi.Service{
					ID:          "service-1",
					Name:        "my-service",
					Description: "my-service-description",
					BrokerName:  "my-broker",
				}))
			})
		})

		When("the service does not exist", func() {
			It("returns error not found", func() {
				Expect(err).To(MatchError("service not found"))
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

				serviceRepo = &Service{
					KubeClient: badKubeClient,
				}
			})

			It("propagates the error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetServices", func() {
		var (
			services []*osbapi.Service
			err      error
		)

		JustBeforeEach(func() {
			services, err = serviceRepo.GetServices("my-broker-1")
		})

		When("services contain owner references to brokers", func() {
			BeforeEach(func() {
				serviceResource := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "my-broker-1",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "broker-uid-1",
						}},
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-1",
						Description: "service-1-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource)).To(Succeed())

				serviceResource2 := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-2",
						Namespace: "default",
						OwnerReferences: []metav1.OwnerReference{{
							Name:       "my-broker-2",
							Kind:       "kind",
							APIVersion: "version",
							UID:        "broker-uid-2",
						}},
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-2",
						Description: "service-2-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource2)).To(Succeed())
			})

			AfterEach(func() {
				deleteServices(kubeClient, "service-1", "service-2")
			})

			It("returns only the services owned by the broker", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(services).To(HaveLen(1))
				Expect(*services[0]).To(MatchFields(IgnoreExtras, Fields{
					"ID":          Equal("service-1"),
					"Name":        Equal("my-service-1"),
					"Description": Equal("service-1-desc"),
					"BrokerName":  Equal("my-broker-1"),
				}))
			})
		})

		When("the service owner reference is not set", func() {
			BeforeEach(func() {
				serviceResource := &v1alpha1.BrokerService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "service-1",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerServiceSpec{
						Name:        "my-service-1",
						Description: "service-1-desc",
					},
				}
				Expect(kubeClient.Create(context.TODO(), serviceResource)).To(Succeed())
			})

			AfterEach(func() {
				deleteServices(kubeClient, "service-1")
			})

			It("successfully returns no services", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(services).To(HaveLen(0))
			})
		})
	})
})

func deleteServices(kubeClient client.Client, serviceNames ...string) {
	for _, s := range serviceNames {
		sToDelete := &v1alpha1.BrokerService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      s,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), sToDelete)).To(Succeed())
	}
}
