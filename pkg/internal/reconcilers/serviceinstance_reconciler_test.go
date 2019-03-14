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
	)

	BeforeEach(func() {
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

	It("updates the service instance status to created", func() {
		Expect(fakeKubeServiceInstanceRepo.UpdateStateCallCount()).To(Equal(1))
		service, newState := fakeKubeServiceInstanceRepo.UpdateStateArgsForCall(0)
		Expect(newState).To(Equal(v1alpha1.ServiceInstanceStateProvisioned))
		Expect(*service).To(Equal(*returnedServiceInstance))
	})

	When("fetching the service instance resource using the kube repo fails", func() {
		BeforeEach(func() {
			fakeKubeServiceInstanceRepo.GetReturns(nil, errors.New("error-getting-serviceinstance"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-serviceinstance"))
		})
	})

	When("the service instane has been deleted", func() {
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

	When("creating the serviceinstance using the broker client fails ", func() {
		BeforeEach(func() {
			fakeBrokerClient.ProvisionInstanceReturns(nil, errors.New("error-provisioning-instance"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-provisioning-instance"))
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

	When("updating the serviceinstance status errors", func() {
		BeforeEach(func() {
			fakeKubeServiceInstanceRepo.UpdateStateReturns(errors.New("error-updating-status"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-updating-status"))
		})
	})
})
