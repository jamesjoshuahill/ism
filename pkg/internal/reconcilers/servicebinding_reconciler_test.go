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

var _ = Describe("ServiceBindingReconciler", func() {
	var (
		reconciler *ServiceBindingReconciler
		err        error

		createBrokerClient         osbapi.CreateFunc
		brokerClientConfiguredWith *osbapi.ClientConfiguration

		returnedServiceBinding *v1alpha1.ServiceBinding
		returnedBroker         *v1alpha1.Broker

		fakeBrokerClient           *reconcilersfakes.FakeBrokerClient
		fakeKubeServiceBindingRepo *reconcilersfakes.FakeKubeServiceBindingRepo
		fakeKubeBrokerRepo         *reconcilersfakes.FakeKubeBrokerRepo
	)

	BeforeEach(func() {
		fakeBrokerClient = &reconcilersfakes.FakeBrokerClient{}
		fakeKubeServiceBindingRepo = &reconcilersfakes.FakeKubeServiceBindingRepo{}
		fakeKubeBrokerRepo = &reconcilersfakes.FakeKubeBrokerRepo{}

		createBrokerClient = func(config *osbapi.ClientConfiguration) (osbapi.Client, error) {
			brokerClientConfiguredWith = config
			return fakeBrokerClient, nil
		}

		returnedServiceBinding = &v1alpha1.ServiceBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-servicebinding-1",
				Namespace: "default",
				UID:       "servicebinding-uid-1",
			},
			Spec: v1alpha1.ServiceBindingSpec{
				Name:       "my-servicebinding-1",
				InstanceID: "instance-1",
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

		fakeKubeServiceBindingRepo.GetReturns(returnedServiceBinding, nil)
		fakeKubeBrokerRepo.GetReturns(returnedBroker, nil)
	})

	JustBeforeEach(func() {
		reconciler = NewServiceBindingReconciler(
			createBrokerClient,
			fakeKubeServiceBindingRepo,
			fakeKubeBrokerRepo,
		)

		_, err = reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{Name: "my-servicebinding-1", Namespace: "default"},
		})
	})

	It("fetches the servicebinding resource using the kube servicebinding repo", func() {
		Expect(err).NotTo(HaveOccurred())

		Expect(fakeKubeServiceBindingRepo.GetCallCount()).To(Equal(1))
		namespacedName := fakeKubeServiceBindingRepo.GetArgsForCall(0)
		Expect(namespacedName).To(Equal(types.NamespacedName{Name: "my-servicebinding-1", Namespace: "default"}))
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

	It("creates the service binding using the broker client", func() {
		Expect(fakeBrokerClient.BindCallCount()).To(Equal(1))
		bindRequest := fakeBrokerClient.BindArgsForCall(0)

		Expect(*bindRequest).To(Equal(osbapi.BindRequest{
			BindingID:         string(returnedServiceBinding.ObjectMeta.UID),
			InstanceID:        returnedServiceBinding.Spec.InstanceID,
			ServiceID:         returnedServiceBinding.Spec.ServiceID,
			PlanID:            returnedServiceBinding.Spec.PlanID,
			AcceptsIncomplete: false,
		}))
	})

	It("updates the service binding status to created", func() {
		Expect(fakeKubeServiceBindingRepo.UpdateStateCallCount()).To(Equal(1))
		service, newState := fakeKubeServiceBindingRepo.UpdateStateArgsForCall(0)
		Expect(newState).To(Equal(v1alpha1.ServiceBindingStateCreated))
		Expect(*service).To(Equal(*returnedServiceBinding))
	})

	When("fetching the service binding resource using the kube repo fails", func() {
		BeforeEach(func() {
			fakeKubeServiceBindingRepo.GetReturns(nil, errors.New("error-getting-servicebinding"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-servicebinding"))
		})
	})

	When("the service binding has been deleted", func() {
		BeforeEach(func() {
			notFoundError := kerrors.NewNotFound(schema.GroupResource{}, "servicebinding")
			fakeKubeServiceBindingRepo.GetReturns(nil, notFoundError)
		})

		It("does not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	When("fetching the servicebinding's broker resource using the kube repo fails", func() {
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

	When("creating the servicebinding using the broker client fails ", func() {
		BeforeEach(func() {
			fakeBrokerClient.BindReturns(nil, errors.New("error-creating-binding"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-creating-binding"))
		})
	})

	When("the servicebinding state reports it is already created ", func() {
		BeforeEach(func() {
			returnedServiceBinding.Status.State = v1alpha1.ServiceBindingStateCreated
		})

		It("doesn't call the broker", func() {
			Expect(fakeBrokerClient.BindCallCount()).To(Equal(0))
		})

		It("doesn't update the status", func() {
			Expect(fakeKubeServiceBindingRepo.UpdateStateCallCount()).To(Equal(0))
		})

		It("still reconciles successfully ", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("updating the servicebinding status errors", func() {
		BeforeEach(func() {
			fakeKubeServiceBindingRepo.UpdateStateReturns(errors.New("error-updating-status"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-updating-status"))
		})
	})
})
