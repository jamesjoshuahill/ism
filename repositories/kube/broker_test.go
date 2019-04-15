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

	"github.com/pivotal-cf/ism/repositories"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/repositories/kube"
)

var _ = Describe("Broker", func() {

	var (
		kubeClient client.Client

		brokerRepo          *Broker
		registrationTimeout time.Duration
	)

	BeforeEach(func() {
		var err error
		kubeClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		registrationTimeout = time.Second

		brokerRepo = &Broker{
			KubeClient:          kubeClient,
			RegistrationTimeout: registrationTimeout,
		}
	})

	Describe("Register", func() {
		var (
			err                  error
			registrationDuration time.Duration
		)

		JustBeforeEach(func() {
			b := &osbapi.Broker{
				Name:     "broker-1",
				URL:      "broker-1-url",
				Username: "broker-1-username",
				Password: "broker-1-password",
			}

			before := time.Now()
			err = brokerRepo.Register(b)
			registrationDuration = time.Since(before)
		})

		When("the status of a broker is set to registered", func() {
			var closeChan chan bool

			BeforeEach(func() {
				closeChan = make(chan bool)
				go simulateRegistrationSuccess(kubeClient, "broker-1", closeChan)
			})

			AfterEach(func() {
				closeChan <- true
				deleteBrokers(kubeClient, "broker-1")
			})

			It("creates a new Broker resource instance", func() {
				Expect(err).NotTo(HaveOccurred())

				key := types.NamespacedName{
					Name:      "broker-1",
					Namespace: "default",
				}

				fetched := &v1alpha1.Broker{}
				Expect(kubeClient.Get(context.TODO(), key, fetched)).To(Succeed())

				Expect(fetched.Spec).To(Equal(v1alpha1.BrokerSpec{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				}))
			})
		})

		When("a broker with the same name already exists", func() {
			var closeChan chan bool

			BeforeEach(func() {
				b := &osbapi.Broker{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				}

				closeChan = make(chan bool)
				go simulateRegistrationSuccess(kubeClient, "broker-1", closeChan)

				// register the broker first, so that the second register errors
				Expect(brokerRepo.Register(b)).To(Succeed())
			})

			AfterEach(func() {
				closeChan <- true
				deleteBrokers(kubeClient, "broker-1")
			})

			It("returns a 'BrokerAlreadyExists' error", func() {
				Expect(err).To(Equal(repositories.ErrBrokerAlreadyExists{BrokerName: "broker-1"}))
			})
		})

		When("the status of a broker is never set to registered", func() {
			It("should eventually timeout", func() {
				Expect(err).To(MatchError(repositories.ErrBrokerRegisterTimeout{BrokerName: "broker-1"}))
			})

			It("times out once the timeout has been reached", func() {
				estimatedExecutionTime := time.Second * 5 // flake prevention!

				Expect(registrationDuration).To(BeNumerically(">", registrationTimeout))
				Expect(registrationDuration).To(BeNumerically("<", registrationTimeout+estimatedExecutionTime))
			})
		})

		When("the status of the broker is set to failed", func() {
			var closeChan chan bool

			BeforeEach(func() {
				closeChan = make(chan bool)
				go simulateRegistrationFailure(kubeClient, "broker-1", "error-message", closeChan)
			})

			AfterEach(func() {
				closeChan <- true
			})

			It("returns the error", func() {
				Expect(err).To(MatchError("error-message"))
			})

			It("deletes the broker resource", func() {
				fetchedBrokerResource := &v1alpha1.Broker{}
				err := kubeClient.Get(context.TODO(), types.NamespacedName{Name: "broker-1", Namespace: "default"}, fetchedBrokerResource)
				Expect(kerrors.IsNotFound(err)).To(BeTrue())
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			brokerResource := &v1alpha1.Broker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-broker",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerSpec{
					Name:     "my-broker",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				},
			}

			Expect(kubeClient.Create(context.TODO(), brokerResource)).To(Succeed())
		})

		It("deletes the Broker resource", func() {
			Expect(brokerRepo.Delete("my-broker")).To(Succeed())

			key := types.NamespacedName{
				Name:      "my-broker",
				Namespace: "default",
			}

			fetched := &v1alpha1.Broker{}
			err := kubeClient.Get(context.TODO(), key, fetched)
			Expect(err).To(HaveOccurred())
			Expect(kerrors.IsNotFound(err)).To(BeTrue())
		})
	})

	Describe("FindByName", func() {
		var (
			broker          *osbapi.Broker
			brokerCreatedAt string
			err             error
		)

		JustBeforeEach(func() {
			broker, err = brokerRepo.FindByName("my-broker")
		})

		When("the broker exists", func() {
			BeforeEach(func() {
				brokerResource := &v1alpha1.Broker{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-broker",
						Namespace: "default",
					},
					Spec: v1alpha1.BrokerSpec{
						Name:     "my-broker",
						URL:      "broker-1-url",
						Username: "broker-1-username",
						Password: "broker-1-password",
					},
				}

				Expect(kubeClient.Create(context.TODO(), brokerResource)).To(Succeed())
				brokerCreatedAt = createdAtForBroker(kubeClient, brokerResource)
			})

			AfterEach(func() {
				deleteBrokers(kubeClient, "my-broker")
			})

			It("returns the broker", func() {
				Expect(err).NotTo(HaveOccurred())

				Expect(*broker).To(Equal(osbapi.Broker{
					Name:      "my-broker",
					URL:       "broker-1-url",
					Username:  "broker-1-username",
					Password:  "broker-1-password",
					CreatedAt: brokerCreatedAt,
				}))
			})
		})

		When("the broker does not exist", func() {
			It("returns an error", func() {
				Expect(err).To(MatchError("Service broker not found"))
			})
		})
	})

	Describe("FindAll", func() {
		var (
			brokers          []*osbapi.Broker
			brokerCreatedAt1 string
			brokerCreatedAt2 string
			err              error
		)

		BeforeEach(func() {
			brokerResource1 := &v1alpha1.Broker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "broker-1",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerSpec{
					Name:     "broker-1",
					URL:      "broker-1-url",
					Username: "broker-1-username",
					Password: "broker-1-password",
				},
			}

			brokerResource2 := &v1alpha1.Broker{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "broker-2",
					Namespace: "default",
				},
				Spec: v1alpha1.BrokerSpec{
					Name:     "broker-2",
					URL:      "broker-2-url",
					Username: "broker-2-username",
					Password: "broker-2-password",
				},
			}

			Expect(kubeClient.Create(context.TODO(), brokerResource1)).To(Succeed())
			Expect(kubeClient.Create(context.TODO(), brokerResource2)).To(Succeed())
			brokerCreatedAt1 = createdAtForBroker(kubeClient, brokerResource1)
			brokerCreatedAt2 = createdAtForBroker(kubeClient, brokerResource2)
		})

		JustBeforeEach(func() {
			brokers, err = brokerRepo.FindAll()
		})

		AfterEach(func() {
			deleteBrokers(kubeClient, "broker-1", "broker-2")
		})

		It("returns all brokers", func() {
			Expect(err).NotTo(HaveOccurred())

			Expect(brokers).To(ConsistOf(
				&osbapi.Broker{
					CreatedAt: brokerCreatedAt1,
					Name:      "broker-1",
					URL:       "broker-1-url",
					Username:  "broker-1-username",
					Password:  "broker-1-password",
				},
				&osbapi.Broker{
					CreatedAt: brokerCreatedAt2,
					Name:      "broker-2",
					URL:       "broker-2-url",
					Username:  "broker-2-username",
					Password:  "broker-2-password",
				},
			))
		})
	})
})

func createdAtForBroker(kubeClient client.Client, brokerResource *v1alpha1.Broker) string {
	b := &v1alpha1.Broker{}
	namespacedName := types.NamespacedName{Name: brokerResource.Name, Namespace: brokerResource.Namespace}

	Expect(kubeClient.Get(context.TODO(), namespacedName, b)).To(Succeed())

	time := b.ObjectMeta.CreationTimestamp.String()
	return time
}

func deleteBrokers(kubeClient client.Client, brokerNames ...string) {
	for _, b := range brokerNames {
		bToDelete := &v1alpha1.Broker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      b,
				Namespace: "default",
			},
		}
		Expect(kubeClient.Delete(context.TODO(), bToDelete)).To(Succeed())
	}
}

func simulateRegistrationSuccess(kubeClient client.Client, brokerName string, done chan bool) {
	simulateRegistration(kubeClient, brokerName, v1alpha1.BrokerStateRegistered, "", done)
}

func simulateRegistrationFailure(kubeClient client.Client, brokerName string, message string, done chan bool) {
	simulateRegistration(kubeClient, brokerName, v1alpha1.BrokerStateRegistrationFailed, message, done)
}

func simulateRegistration(kubeClient client.Client, brokerName string, state v1alpha1.BrokerState, message string, done chan bool) {
	for {
		select {
		case <-done:
			return //exit func
		default:
			key := types.NamespacedName{
				Name:      brokerName,
				Namespace: "default",
			}
			broker := &v1alpha1.Broker{}
			err := kubeClient.Get(context.TODO(), key, broker)
			if err != nil {
				break //loop again
			}

			if broker.Status.State != state || broker.Status.Message != message {
				broker.Status.State = state
				broker.Status.Message = message
				Expect(kubeClient.Status().Update(context.TODO(), broker)).To(Succeed())
			}
		}
	}
}
