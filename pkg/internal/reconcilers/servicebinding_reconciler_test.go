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
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var _ = Describe("ServiceBindingReconciler", func() {
	var (
		reconciler *ServiceBindingReconciler
		err        error

		log                        logr.Logger
		createBrokerClient         osbapi.CreateFunc
		brokerClientConfiguredWith *osbapi.ClientConfiguration

		returnedServiceBinding *v1alpha1.ServiceBinding
		returnedBroker         *v1alpha1.Broker

		fakeBrokerClient           *reconcilersfakes.FakeBrokerClient
		fakeKubeServiceBindingRepo *reconcilersfakes.FakeKubeServiceBindingRepo
		fakeKubeBrokerRepo         *reconcilersfakes.FakeKubeBrokerRepo
		fakeKubeSecretRepo         *reconcilersfakes.FakeKubeSecretRepo
	)

	BeforeEach(func() {
		fakeBrokerClient = &reconcilersfakes.FakeBrokerClient{}
		fakeKubeServiceBindingRepo = &reconcilersfakes.FakeKubeServiceBindingRepo{}
		fakeKubeBrokerRepo = &reconcilersfakes.FakeKubeBrokerRepo{}
		fakeKubeSecretRepo = &reconcilersfakes.FakeKubeSecretRepo{}
		log = logf.Log.WithName("test-logger")

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

		fakeBrokerClient.BindReturns(&osbapi.BindResponse{Credentials: map[string]interface{}{"password": "my-secret"}}, nil)
		fakeKubeServiceBindingRepo.GetReturns(returnedServiceBinding, nil)
		fakeKubeBrokerRepo.GetReturns(returnedBroker, nil)
		fakeKubeSecretRepo.CreateReturns(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "some-secret"}}, nil)
	})

	JustBeforeEach(func() {
		reconciler = NewServiceBindingReconciler(
			log,
			createBrokerClient,
			fakeKubeServiceBindingRepo,
			fakeKubeBrokerRepo,
			fakeKubeSecretRepo,
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

	When("creating the service binding", func() {
		It("updates the status to creating", func() {
			binding, newState := fakeKubeServiceBindingRepo.UpdateStateArgsForCall(0)
			Expect(newState).To(Equal(v1alpha1.ServiceBindingStateCreating))
			Expect(*binding).To(Equal(*returnedServiceBinding))
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

		It("creates a secret with the broker returned credentials", func() {
			Expect(fakeKubeSecretRepo.CreateCallCount()).To(Equal(1))
			bindingArg, credArg := fakeKubeSecretRepo.CreateArgsForCall(0)

			Expect(bindingArg).To(Equal(returnedServiceBinding))
			Expect(credArg).To(Equal(map[string]interface{}{"password": "my-secret"}))
		})

		It("updates the service binding status secretRef to the secret created", func() {
			binding, _ := fakeKubeServiceBindingRepo.UpdateStateArgsForCall(1)
			Expect(binding.Status.SecretRef.Name).To(Equal("some-secret"))
		})

		It("updates the service binding status state to created", func() {
			binding, newState := fakeKubeServiceBindingRepo.UpdateStateArgsForCall(1)
			Expect(newState).To(Equal(v1alpha1.ServiceBindingStateCreated))
			Expect(*binding).To(Equal(*returnedServiceBinding))
		})

		It("adds a finalizer", func() {
			Expect(fakeKubeServiceBindingRepo.UpdateCallCount()).To(Equal(1))
			binding := fakeKubeServiceBindingRepo.UpdateArgsForCall(0)
			Expect(binding.ObjectMeta.Finalizers).To(HaveLen(1))
			Expect(binding.ObjectMeta.Finalizers[0]).To(Equal("finalizer.servicebinding.osbapi.ism.io"))
		})

		When("creating the servicebinding using the broker client fails ", func() {
			BeforeEach(func() {
				fakeBrokerClient.BindReturns(nil, errors.New("error-creating-binding"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-creating-binding"))
			})
		})

		When("creating the secret using the kube secret repo fails", func() {
			BeforeEach(func() {
				fakeKubeSecretRepo.CreateReturns(nil, errors.New("error-creating-secret"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-creating-secret"))
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

	When("the service binding has been marked for deletion", func() {
		BeforeEach(func() {
			time := metav1.Now()
			returnedServiceBinding.ObjectMeta.DeletionTimestamp = &time

			returnedServiceBinding.Status.SecretRef = corev1.LocalObjectReference{Name: "my-secret"}
			returnedServiceBinding.Status.State = v1alpha1.ServiceBindingStateCreated
		})

		It("updates the status to deleting", func() {
			Expect(fakeKubeServiceBindingRepo.UpdateStateCallCount()).To(Equal(1))
			binding, newState := fakeKubeServiceBindingRepo.UpdateStateArgsForCall(0)
			Expect(newState).To(Equal(v1alpha1.ServiceBindingStateDeleting))
			Expect(*binding).To(Equal(*returnedServiceBinding))
		})

		It("deletes the service binding using the broker client", func() {
			Expect(fakeBrokerClient.UnbindCallCount()).To(Equal(1))
			unbindRequest := fakeBrokerClient.UnbindArgsForCall(0)

			Expect(*unbindRequest).To(Equal(osbapi.UnbindRequest{
				BindingID:         string(returnedServiceBinding.ObjectMeta.UID),
				InstanceID:        returnedServiceBinding.Spec.InstanceID,
				ServiceID:         returnedServiceBinding.Spec.ServiceID,
				PlanID:            returnedServiceBinding.Spec.PlanID,
				AcceptsIncomplete: false,
			}))
		})

		It("deletes the secret associated with the service binding", func() {
			Expect(fakeKubeSecretRepo.DeleteCallCount()).To(Equal(1))
			deleteReq := fakeKubeSecretRepo.DeleteArgsForCall(0)

			Expect(deleteReq.Name).To(Equal(returnedServiceBinding.Status.SecretRef.Name))
			Expect(deleteReq.Namespace).To(Equal(returnedServiceBinding.Namespace))
		})

		It("removes the finalizer", func() {
			Expect(fakeKubeServiceBindingRepo.UpdateCallCount()).To(Equal(1))
			binding := fakeKubeServiceBindingRepo.UpdateArgsForCall(0)
			Expect(binding.ObjectMeta.Finalizers).To(HaveLen(0))
		})

		When("deleting the servicebinding using the broker client fails ", func() {
			BeforeEach(func() {
				fakeBrokerClient.UnbindReturns(nil, errors.New("error-deleting-binding"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-deleting-binding"))
			})
		})

		When("deleting the secret using the kube secret repo fails", func() {
			BeforeEach(func() {
				fakeKubeSecretRepo.DeleteReturns(errors.New("error-deleting-secret"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-deleting-secret"))
			})
		})

		When("removing the finalizers errors when updating the servicebinding", func() {
			BeforeEach(func() {
				fakeKubeServiceBindingRepo.UpdateReturns(errors.New("error-updating"))
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-updating"))
			})
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

	When("fetching the service binding resource using the kube repo fails", func() {
		BeforeEach(func() {
			fakeKubeServiceBindingRepo.GetReturns(nil, errors.New("error-getting-servicebinding"))
		})

		It("returns the error", func() {
			Expect(err).To(MatchError("error-getting-servicebinding"))
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

})
