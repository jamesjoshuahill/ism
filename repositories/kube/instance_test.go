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
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"github.com/pivotal-cf/ism/repositories"
	. "github.com/pivotal-cf/ism/repositories/kube"
)

var _ = Describe("Instance", func() {

	var (
		kubeClient client.Client

		instanceRepo *Instance
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		instanceRepo = &Instance{
			KubeClient: kubeClient,
		}
	})

	Describe("GetInstances", func() {
		var (
			instances          []*osbapi.Instance
			instanceID1        string
			instanceID2        string
			instanceCreatedAt1 string
			instanceCreatedAt2 string
			err                error
		)

		BeforeEach(func() {
			instanceResource1 := &v1alpha1.ServiceInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-1",
					Namespace: "default",
				},
				Spec: v1alpha1.ServiceInstanceSpec{
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "my-broker-1",
				},
				Status: v1alpha1.ServiceInstanceStatus{
					State: v1alpha1.ServiceInstanceStateProvisioned,
				},
			}

			instanceResource2 := &v1alpha1.ServiceInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-2",
					Namespace: "default",
				},
				Spec: v1alpha1.ServiceInstanceSpec{
					Name:       "instance-2",
					PlanID:     "plan-2",
					ServiceID:  "service-2",
					BrokerName: "my-broker-2",
				},
			}

			Expect(kubeClient.Create(context.TODO(), instanceResource1)).To(Succeed())
			Expect(kubeClient.Status().Update(context.TODO(), instanceResource1)).To(Succeed())

			Expect(kubeClient.Create(context.TODO(), instanceResource2)).To(Succeed())
			instanceCreatedAt1 = createdAtForInstance(kubeClient, instanceResource1)
			instanceCreatedAt2 = createdAtForInstance(kubeClient, instanceResource2)
			instanceID1 = idForInstance(kubeClient, instanceResource1)
			instanceID2 = idForInstance(kubeClient, instanceResource2)
		})

		JustBeforeEach(func() {
			instances, err = instanceRepo.GetInstances()
		})

		AfterEach(func() {
			deleteInstances(kubeClient, "instance-1", "instance-2")
		})

		It("returns all instances", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(instances).To(ConsistOf(
				&osbapi.Instance{
					ID:         instanceID1,
					CreatedAt:  instanceCreatedAt1,
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					Status:     "created",
					BrokerName: "my-broker-1",
				},
				&osbapi.Instance{
					ID:         instanceID2,
					CreatedAt:  instanceCreatedAt2,
					Name:       "instance-2",
					PlanID:     "plan-2",
					ServiceID:  "service-2",
					Status:     "creating",
					BrokerName: "my-broker-2",
				},
			))
		})
	})

	Describe("GetInstanceByName", func() {
		var (
			instance          *osbapi.Instance
			instanceID        string
			instanceCreatedAt string
			err               error
		)

		JustBeforeEach(func() {
			instance, err = instanceRepo.GetInstanceByName("instance")
		})

		When("the instance exists", func() {
			BeforeEach(func() {
				instanceResource1 := &v1alpha1.ServiceInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance",
						Namespace: "default",
					},
					Spec: v1alpha1.ServiceInstanceSpec{
						Name:       "instance",
						PlanID:     "plan-1",
						ServiceID:  "service-1",
						BrokerName: "my-broker-1",
					},
				}
				Expect(kubeClient.Create(context.TODO(), instanceResource1)).To(Succeed())
				instanceCreatedAt = createdAtForInstance(kubeClient, instanceResource1)
				instanceID = idForInstance(kubeClient, instanceResource1)
			})

			AfterEach(func() {
				deleteInstances(kubeClient, "instance")
			})

			It("returns the instance", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(instance).To(Equal(
					&osbapi.Instance{
						ID:         instanceID,
						CreatedAt:  instanceCreatedAt,
						Name:       "instance",
						PlanID:     "plan-1",
						ServiceID:  "service-1",
						Status:     "creating",
						BrokerName: "my-broker-1",
					},
				))
			})
		})

		When("the instance does not exist", func() {
			It("returns a not found error", func() {
				Expect(err).To(MatchError("instance not found"))
			})
		})
	})

	Describe("GetInstancesForBroker", func() {
		var (
			instances          []*osbapi.Instance
			instanceID1        string
			instanceCreatedAt1 string
			err                error
		)

		BeforeEach(func() {
			instanceResource1 := &v1alpha1.ServiceInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-1",
					Namespace: "default",
				},
				Spec: v1alpha1.ServiceInstanceSpec{
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "my-broker-1",
				},
				Status: v1alpha1.ServiceInstanceStatus{
					State: v1alpha1.ServiceInstanceStateProvisioned,
				},
			}

			instanceResource2 := &v1alpha1.ServiceInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance-2",
					Namespace: "default",
				},
				Spec: v1alpha1.ServiceInstanceSpec{
					Name:       "instance-2",
					PlanID:     "plan-2",
					ServiceID:  "service-2",
					BrokerName: "my-broker-2",
				},
			}

			Expect(kubeClient.Create(context.TODO(), instanceResource1)).To(Succeed())
			Expect(kubeClient.Status().Update(context.TODO(), instanceResource1)).To(Succeed())

			Expect(kubeClient.Create(context.TODO(), instanceResource2)).To(Succeed())
			instanceCreatedAt1 = createdAtForInstance(kubeClient, instanceResource1)
			instanceID1 = idForInstance(kubeClient, instanceResource1)
		})

		JustBeforeEach(func() {
			instances, err = instanceRepo.GetInstancesForBroker("my-broker-1")
		})

		AfterEach(func() {
			deleteInstances(kubeClient, "instance-1", "instance-2")
		})

		It("returns all instances for a specific broker", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(instances).To(ConsistOf(
				&osbapi.Instance{
					ID:         instanceID1,
					CreatedAt:  instanceCreatedAt1,
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					Status:     "created",
					BrokerName: "my-broker-1",
				},
			))
		})
	})

	Describe("GetInstanceByID", func() {
		var (
			instance          *osbapi.Instance
			instanceCreatedAt string
			err               error
			instanceID        string
		)

		JustBeforeEach(func() {
			instance, err = instanceRepo.GetInstanceByID(instanceID)
		})

		When("the instance exists", func() {
			BeforeEach(func() {
				instanceResource := &v1alpha1.ServiceInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "instance",
						Namespace: "default",
						UID:       "instance-1",
					},
					Spec: v1alpha1.ServiceInstanceSpec{
						Name:       "instance",
						PlanID:     "plan-1",
						ServiceID:  "service-1",
						BrokerName: "my-broker-1",
					},
				}
				Expect(kubeClient.Create(context.TODO(), instanceResource)).To(Succeed())
				instanceCreatedAt = createdAtForInstance(kubeClient, instanceResource)
				instanceID = idForInstance(kubeClient, instanceResource)
			})

			AfterEach(func() {
				deleteInstances(kubeClient, "instance")
			})

			It("returns the instance", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(instance).To(Equal(
					&osbapi.Instance{
						ID:         instanceID,
						Name:       "instance",
						PlanID:     "plan-1",
						ServiceID:  "service-1",
						Status:     "creating",
						BrokerName: "my-broker-1",
						CreatedAt:  instanceCreatedAt,
					},
				))
			})
		})

		When("the instance does not exist", func() {
			It("returns a not found error", func() {
				Expect(err).To(MatchError("instance not found"))
			})
		})
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			b := &osbapi.Instance{
				Name:       "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			}

			err = instanceRepo.Create(b)
		})

		AfterEach(func() {
			deleteInstances(kubeClient, "instance-1")
		})

		It("creates a new ServiceInstance resource", func() {
			Expect(err).NotTo(HaveOccurred())

			key := types.NamespacedName{
				Name:      "instance-1",
				Namespace: "default",
			}

			fetched := &v1alpha1.ServiceInstance{}
			Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

			Expect(fetched.Spec).To(Equal(v1alpha1.ServiceInstanceSpec{
				Name:       "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			}))
		})

		When("creating an Instance with a duplicate name", func() {
			BeforeEach(func() {
				// create the instance first, so that the second create errors
				b := &osbapi.Instance{
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "broker-1",
				}

				Expect(instanceRepo.Create(b)).To(Succeed())
			})

			It("returns a 'ErrInstanceAlreadyExists' error", func() {
				Expect(err).To(Equal(repositories.ErrInstanceAlreadyExists{InstanceName: "instance-1"}))
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			instanceResource := &v1alpha1.ServiceInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-instance",
					Namespace: "default",
				},
				Spec: v1alpha1.ServiceInstanceSpec{
					Name:       "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "broker-1",
				},
			}

			Expect(kubeClient.Create(context.TODO(), instanceResource)).To(Succeed())
		})

		It("deletes the ServiceInstance resource", func() {
			Expect(instanceRepo.Delete("my-instance")).To(Succeed())

			key := types.NamespacedName{
				Name:      "my-instance",
				Namespace: "default",
			}

			fetched := &v1alpha1.ServiceInstance{}
			err := kubeClient.Get(context.TODO(), key, fetched)
			Expect(err).To(HaveOccurred())
			Expect(kerrors.IsNotFound(err)).To(BeTrue())
		})
	})
})

func createdAtForInstance(kubeClient client.Client, instanceResource *v1alpha1.ServiceInstance) string {
	i := &v1alpha1.ServiceInstance{}
	namespacedName := types.NamespacedName{Name: instanceResource.Name, Namespace: instanceResource.Namespace}

	Expect(kubeClient.Get(context.TODO(), namespacedName, i)).To(Succeed())

	time := i.ObjectMeta.CreationTimestamp.String()
	return time
}

func idForInstance(kubeClient client.Client, instanceResource *v1alpha1.ServiceInstance) string {
	i := &v1alpha1.ServiceInstance{}
	namespacedName := types.NamespacedName{Name: instanceResource.Name, Namespace: instanceResource.Namespace}

	Expect(kubeClient.Get(context.TODO(), namespacedName, i)).To(Succeed())

	return string(i.ObjectMeta.UID)
}

func deleteInstances(kubeClient client.Client, instanceNames ...string) {
	for _, name := range instanceNames {
		iToDelete := &v1alpha1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), iToDelete)).To(Succeed())
	}
}
