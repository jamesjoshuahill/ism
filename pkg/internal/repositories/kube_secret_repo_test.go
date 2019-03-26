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
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	. "github.com/pivotal-cf/ism/pkg/internal/repositories"
)

var _ = Describe("KubeSecretRepo", func() {
	var (
		repo           *KubeSecretRepo
		binding        *v1alpha1.ServiceBinding
		returnedSecret *corev1.Secret
	)

	BeforeEach(func() {
		binding = &v1alpha1.ServiceBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-binding",
				Namespace: "default",
			},
		}

		err := kubeClient.Create(context.Background(), binding)
		Expect(err).NotTo(HaveOccurred())
		Expect(kubeClient.Get(context.TODO(), types.NamespacedName{Name: binding.Name, Namespace: binding.Namespace}, binding)).To(Succeed())

		repo = NewKubeSecretRepo(kubeClient)
	})

	AfterEach(func() {
		Expect(kubeClient.Delete(context.Background(), binding)).To(Succeed())
	})

	Describe("Create", func() {
		When("the credential can be serialized", func() {
			BeforeEach(func() {
				var err error
				returnedSecret, err = repo.Create(binding, map[string]interface{}{"username": "admin"})
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(kubeClient.Delete(context.Background(), returnedSecret)).To(Succeed())
			})

			It("creates the secret with a 'ism-cred-' prefix", func() {
				createdSecret := corev1.Secret{}

				Expect(kubeClient.Get(context.TODO(), types.NamespacedName{Name: "ism-cred-my-binding", Namespace: binding.Namespace}, &createdSecret)).To(Succeed())
				Expect(createdSecret.Name).To(Equal("ism-cred-my-binding"))
				Expect(createdSecret.Namespace).To(Equal(binding.Namespace))

				rawCreds, err := json.Marshal(map[string]string{"username": "admin"})
				Expect(err).NotTo(HaveOccurred())
				Expect(createdSecret.Data["credentials"]).To(Equal(rawCreds))

				Expect(returnedSecret.Name).To(Equal("ism-cred-my-binding"))
			})

			It("sets the owner references on the binding", func() {
				Expect(returnedSecret.ObjectMeta.OwnerReferences).To(HaveLen(1))
				Expect(returnedSecret.ObjectMeta.OwnerReferences[0].UID).To(Equal(binding.ObjectMeta.UID))
			})
		})

		When("the credential cannot be serialized", func() {
			It("returns an error", func() {
				_, err := repo.Create(binding, map[string]interface{}{"username": make(chan int)})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
