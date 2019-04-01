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

package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const secretPrefix = "ism-cred-"

type KubeSecretRepo struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewKubeSecretRepo(client client.Client) *KubeSecretRepo {
	return &KubeSecretRepo{
		client: client,
		scheme: scheme.Scheme,
	}
}

func (repo *KubeSecretRepo) Create(binding *v1alpha1.ServiceBinding, creds map[string]interface{}) (*corev1.Secret, error) {
	encodedCreds, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretPrefix + binding.Name,
			Namespace: binding.Namespace,
		},
		Data: map[string][]byte{"credentials": encodedCreds},
	}

	if err := controllerutil.SetControllerReference(binding, secret, repo.scheme); err != nil {
		return nil, err
	}

	if err := repo.client.Create(context.TODO(), secret); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return secret, nil
}

func (repo *KubeSecretRepo) Delete(secret *corev1.Secret) error {
	return repo.client.Delete(context.TODO(), secret)
}
