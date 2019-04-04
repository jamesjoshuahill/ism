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

package reconcilers_test

import (
	"errors"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/reconcilers"
	"github.com/pivotal-cf/ism/pkg/internal/reconcilers/reconcilersfakes"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var _ = Describe("ServiceInstanceReconciler", func() {
	var (
		reconciler *ServiceInstanceReconciler
		err        error

		createBrokerClient         osbapi.CreateFunc
		brokerClientConfiguredWith *osbapi.ClientConfiguration

		returnedServiceInstance *v1alpha1.ServiceInstance
		returnedBroker          *v1alpha1.Broker

		fakeBrokerClient            *reconcilersfakes.FakeBrokerClient
		fakeKubeServiceInstanceRepo *reconcilersfakes.FakeKubeServiceInstanceRepo
		fakeKubeBrokerRepo          *reconcilersfakes.FakeKubeBrokerRepo

		logger logr.Logger
	)

	BeforeEach(func() {
		logger = logf.Log.WithName("dummy-logger")
		fakeBrokerClient = &reconcilersfakes.FakeBrokerClient{}
		fakeKubeServiceInstanceRepo = &reconcilersfakes.FakeKubeServiceInstanceRepo{}
		fakeKubeBrokerRepo = &reconcilersfakes.FakeKubeBrokerRepo{}

		createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
			brokerClientConfiguredWith = config
			return fakeBrokerClient, nil
		}

		returnedServiceInstance = &v1alpha1.ServiceInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-serviceinstance-1",
				Namespace: "default",
				UID:       "serviceinstance-uid-1",
			},
			Spec: v1alpha1.ServiceInstanceSpec{
				Name:       "my-serviceinstance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			},
		}
		returnedBroker = &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "broker-1",
				Namespace: "default",
			},
			Spec: v1alpha1.BrokerSpec{
				Name:     "broker-1",
				URL:      "broker-url",
				Username: "broker-username",
				Password: "broker-password",
			},
		}

		fakeKubeServiceInstanceRepo.GetReturns(returnedServiceInstance, nil)
		fakeKubeBrokerRepo.GetReturns(returnedBroker, nil)
	})

	JustBeforeEach(func() {
		reconciler = NewServiceInstanceReconciler(
			logger,
			createBrokerClient,
			fakeKubeServiceInstanceRepo,
			fakeKubeBrokerRepo,
		)

		_, err = reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{Name: "my-serviceinstance-1", Namespace: "default"},
		})
	})

	It("fetches the serviceinstance resource using the kube serviceinstance repo", func() {
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeKubeServiceInstanceRepo.GetCallCount()).To(Equal(1))
		namespacedName := fakeKubeServiceInstanceRepo.GetArgsForCall(0)
		Expect(namespacedName).To(Equal(types.NamespacedName{Name: "my-serviceinstance-1", Namespace: "default"}))
	})

	It("configures the broker client with correct options", func() {
		Expect(*brokerClientConfiguredWith).To(Equal(osbapi.ClientConfiguration{
			Name:                "broker-1",
			URL:                 "broker-url",
			APIVersion:          osbapi.LatestAPIVersion(),
			TimeoutSeconds:      60,
			EnableAlphaFeatures: false,
			AuthConfig: &osbapi.AuthConfig{
				BasicAuthConfig: &osbapi.BasicAuthConfig{
					Username: "broker-username",
					Password: "broker-password",
				},
			},
		}))
	})

	When("provisioning the service instance", func() {
		It("updates the status to provisioning", func() {
			instance, newState := fakeKubeServiceInstanceRepo.UpdateStateArgsForCall(0)
			Expect(newState).To(Equal(v1alpha1.ServiceInstanceStateProvisioning))
			Expect(*instance).To(Equal(*returnedServiceInstance))
		})

		It("creates the service instance using the broker client", func() {
			Expect(fakeBrokerClient.ProvisionInstanceCallCount()).To(Equal(1))
			provisionRequest := fakeBrokerClient.ProvisionInstanceArgsForCall(0)

			Expect(*provisionRequest).To(Equal(osbapi.ProvisionRequest{
				InstanceID:        string(returnedServiceInstance.ObjectMeta.UID),
				AcceptsIncomplete: false,
				ServiceID:         returnedServiceInstance.Spec.ServiceID,
				PlanID:            returnedServiceInstance.Spec.PlanID,
				OrganizationGUID:  returnedServiceInstance.ObjectMeta.Namespace,
				SpaceGUID:         returnedServiceInstance.ObjectMeta.Namespace,
			}))
		})

		It("adds a finalizer", func() {
			Expect(fakeKubeServiceInstanceRepo.UpdateCallCount()).To(Equal(1))
			instance := fakeKubeServiceInstanceRepo.UpdateArgsForCall(0)
			Expect(instance.ObjectMeta.Finalizers).To(HaveLen(1))
			Expect(instance.ObjectMeta.Finalizers[0]).To(Equal("finalizer.serviceinstance.osbapi.ism.io"))
		})

		It("updates the status to provisioned", func() {
			service, newState := fakeKubeServiceInstanceRepo.UpdateStateArgsForCall(1)
			Expect(newState).To(Equal(v1alpha1.ServiceInstanceStateProvisioned))
			Expect(*service).To(Equal(*returnedServiceInstance))
		})

		When("updating the serviceinstance status to provisioning errors", func() {
			BeforeEach(func() {
				fakeKubeServiceInstanceRepo.UpdateStateReturnsOnCall(0, errors.New("error-updating-status"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating-status"))
			})
		})

		When("creating the serviceinstance using the broker client fails ", func() {
			BeforeEach(func() {
				fakeBrokerClient.ProvisionInstanceReturns(nil, errors.New("error-provisioning-instance"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-provisioning-instance"))
			})
		})

		When("updating the serviceinstance finalizer errors", func() {
			BeforeEach(func() {
				fakeKubeServiceInstanceRepo.UpdateReturns(errors.New("error-updating"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating"))
			})
		})

		When("updating the serviceinstance status to provisioned errors", func() {
			BeforeEach(func() {
				fakeKubeServiceInstanceRepo.UpdateStateReturnsOnCall(1, errors.New("error-updating-status"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating-status"))
			})
		})
	})

	When("the service instance has been marked for deletion", func() {
		BeforeEach(func() {
			time := metav1.Now()
			returnedServiceInstance.ObjectMeta.DeletionTimestamp = &time

			returnedServiceInstance.Status.State = v1alpha1.ServiceInstanceStateProvisioned
		})

		It("updates the status to deprovisioning", func() {
			Expect(fakeKubeServiceInstanceRepo.UpdateStateCallCount()).To(Equal(1))
			instance, newState := fakeKubeServiceInstanceRepo.UpdateStateArgsForCall(0)
			Expect(newState).To(Equal(v1alpha1.ServiceInstanceStateDeprovisioning))
			Expect(*instance).To(Equal(*returnedServiceInstance))
		})

		It("deletes the service instance using the broker client", func() {
			Expect(fakeBrokerClient.DeprovisionInstanceCallCount()).To(Equal(1))
			deprovisionRequest := fakeBrokerClient.DeprovisionInstanceArgsForCall(0)

			Expect(*deprovisionRequest).To(Equal(osbapi.DeprovisionRequest{
				InstanceID:        string(returnedServiceInstance.ObjectMeta.UID),
				ServiceID:         returnedServiceInstance.Spec.ServiceID,
				PlanID:            returnedServiceInstance.Spec.PlanID,
				AcceptsIncomplete: false,
			}))
		})

		It("removes the finalizer", func() {
			Expect(fakeKubeServiceInstanceRepo.UpdateCallCount()).To(Equal(1))
			instance := fakeKubeServiceInstanceRepo.UpdateArgsForCall(0)
			Expect(instance.ObjectMeta.Finalizers).To(HaveLen(0))
		})

		When("updating the status to deprovisioning errors", func() {
			BeforeEach(func() {
				fakeKubeServiceInstanceRepo.UpdateStateReturnsOnCall(0, errors.New("error-updating-status"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating-status"))
			})
		})

		When("deprovisioning the serviceinstance using the broker client fails ", func() {
			BeforeEach(func() {
				fakeBrokerClient.DeprovisionInstanceReturns(nil, errors.New("error-deprovisioning-instance"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-deprovisioning-instance"))
			})
		})

		When("updating the serviceinstance finalizer errors", func() {
			BeforeEach(func() {
				fakeKubeServiceInstanceRepo.UpdateReturns(errors.New("error-updating"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating"))
			})
		})
	})

	When("fetching the service instance resource using the kube repo fails", func() {
		BeforeEach(func() {
			fakeKubeServiceInstanceRepo.GetReturns(nil, errors.New("error-getting-serviceinstance"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-serviceinstance"))
		})
	})

	When("the service instance has been deleted", func() {
		BeforeEach(func() {
			notFoundError := kerrors.NewNotFound(schema.GroupResource{}, "serviceinstance")
			fakeKubeServiceInstanceRepo.GetReturns(nil, notFoundError)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("fetching the serviceinstance's broker resource using the kube repo fails", func() {
		BeforeEach(func() {
			fakeKubeBrokerRepo.GetReturns(nil, errors.New("error-getting-broker"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-broker"))
		})
	})

	When("configuring the broker client fails", func() {
		BeforeEach(func() {
			createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
				return nil, errors.New("error-configuring-broker-client")
			}
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-configuring-broker-client"))
		})
	})

	When("the serviceinstance state reports it is already provisioned", func() {
		BeforeEach(func() {
			returnedServiceInstance.Status.State = v1alpha1.ServiceInstanceStateProvisioned
		})

		It("doesn't call the broker", func() {
			Expect(fakeBrokerClient.ProvisionInstanceCallCount()).To(Equal(0))
		})

		It("doesn't update the status", func() {
			Expect(fakeKubeServiceInstanceRepo.UpdateStateCallCount()).To(Equal(0))
		})

		It("still reconciles successfully ", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
