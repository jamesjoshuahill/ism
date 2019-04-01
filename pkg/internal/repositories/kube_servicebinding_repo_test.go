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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeServiceBindingRepo", func() {
	var (
		repo                   *KubeServiceBindingRepo
		existingServiceBinding *v1alpha1.ServiceBinding
		resource               types.NamespacedName
	)

	BeforeEach(func() {
		resource = types.NamespacedName{Name: "my-servicebinding-1", Namespace: "default"}

		existingServiceBinding = &v1alpha1.ServiceBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resource.Name,
				Namespace: resource.Namespace,
			},
			Spec: v1alpha1.ServiceBindingSpec{
				Name:       "my-servicebinding-1",
				InstanceID: "instance-1",
				PlanID:     "plan-1",
				ServiceID:  "service-1",
				BrokerName: "broker-1",
			},
		}

		repo = NewKubeServiceBindingRepo(kubeClient)
	})

	AfterEach(func() {
		kubeClient.Delete(context.Background(), existingServiceBinding)
	})

	Describe("Get", func() {
		When("the servicebinding exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingServiceBinding)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns servicebinding", func() {
				servicebinding, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(servicebinding).To(Equal(existingServiceBinding))
			})
		})

		When("the servicebinding doesn't exist", func() {
			It("returns an error", func() {
				_, err := repo.Get(resource)
				Expect(err).To(MatchError("servicebindings.osbapi.ism.io \"my-servicebinding-1\" not found"))
			})
		})
	})

	Describe("Update", func() {
		When("the servicebinding exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingServiceBinding)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates the servicebinding", func() {
				existingServiceBinding.ObjectMeta.SetLabels(map[string]string{"label": "one"})

				Expect(repo.Update(existingServiceBinding)).To(Succeed())

				remoteModifiedBinding := &v1alpha1.ServiceBinding{}

				Expect(kubeClient.Get(
					context.TODO(),
					types.NamespacedName{
						Name:      existingServiceBinding.Name,
						Namespace: existingServiceBinding.Namespace},
					remoteModifiedBinding)).To(Succeed())

				Expect(remoteModifiedBinding.ObjectMeta.GetLabels()).To(Equal(map[string]string{"label": "one"}))
			})
		})

		When("the servicebinding doesn't exist", func() {
			It("returns an error", func() {
				err := repo.Update(existingServiceBinding)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UpdateStatus", func() {
		When("the servicebinding exists", func() {
			BeforeEach(func() {
				err := kubeClient.Create(context.Background(), existingServiceBinding)
				Expect(err).NotTo(HaveOccurred())
			})

			It("updates status", func() {
				newState := v1alpha1.ServiceBindingStateCreated
				Expect(existingServiceBinding.Status.State).NotTo(Equal(newState))

				err := repo.UpdateState(existingServiceBinding, newState)
				Expect(err).NotTo(HaveOccurred())

				updatedServiceBinding, err := repo.Get(resource)
				Expect(err).NotTo(HaveOccurred())

				Expect(updatedServiceBinding.Status.State).To(Equal(newState))
				Expect(existingServiceBinding.Status.State).To(Equal(newState))
			})
		})

		When("the servicebinding doesn't exist", func() {
			It("returns an error", func() {
				newState := v1alpha1.ServiceBindingStateCreated
				err := repo.UpdateState(existingServiceBinding, newState)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
