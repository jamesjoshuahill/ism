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

type KubePlanRepo struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewKubePlanRepo(client client.Client) *KubePlanRepo {
	return &KubePlanRepo{
		client: client,
		scheme: scheme.Scheme,
	}
}

func (repo *KubePlanRepo) Create(service *v1alpha1.BrokerService, catalogPlan osbapi.Plan) error {
	plan := &v1alpha1.BrokerServicePlan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      catalogPlan.ID,
			Namespace: service.Namespace,
		},
		Spec: v1alpha1.BrokerServicePlanSpec{
			Name: catalogPlan.Name,
		},
	}

	if err := controllerutil.SetControllerReference(service, plan, repo.scheme); err != nil {
		return err
	}

	return repo.client.Create(context.TODO(), plan)
}
