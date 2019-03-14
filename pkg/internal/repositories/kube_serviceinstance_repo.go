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
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// var ctx = context.TODO()

type KubeServiceInstanceRepo struct {
	client client.Client
}

func NewKubeServiceInstanceRepo(client client.Client) *KubeServiceInstanceRepo {
	return &KubeServiceInstanceRepo{
		client: client,
	}
}

func (repo *KubeServiceInstanceRepo) Get(resource types.NamespacedName) (*v1alpha1.ServiceInstance, error) {
	serviceInstance := &v1alpha1.ServiceInstance{}

	err := repo.client.Get(ctx, resource, serviceInstance)
	if err != nil {
		return nil, err
	}

	return serviceInstance, nil
}

func (repo *KubeServiceInstanceRepo) UpdateState(serviceInstance *v1alpha1.ServiceInstance, newState v1alpha1.ServiceInstanceState) error {
	serviceInstance.Status.State = newState

	return repo.client.Status().Update(ctx, serviceInstance)
}
