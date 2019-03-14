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

	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	osbapi "github.com/pmorie/go-open-service-broker-client/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type KubeServiceRepo struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewKubeServiceRepo(client client.Client) *KubeServiceRepo {
	return &KubeServiceRepo{
		client: client,
		scheme: scheme.Scheme,
	}
}

func (repo *KubeServiceRepo) Create(broker *v1alpha1.Broker, catalogService osbapi.Service) (*v1alpha1.BrokerService, error) {
	service := &v1alpha1.BrokerService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      catalogService.ID,
			Namespace: broker.Namespace,
		},
		Spec: v1alpha1.BrokerServiceSpec{
			BrokerID:    broker.Name,
			Name:        catalogService.Name,
			Description: catalogService.Description,
		},
	}

	if err := controllerutil.SetControllerReference(broker, service, repo.scheme); err != nil {
		return nil, err
	}

	if err := repo.client.Create(context.TODO(), service); err != nil {
		return nil, err
	}

	return service, nil
}
