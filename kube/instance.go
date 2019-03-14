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

package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Instance struct {
	KubeClient client.Client
}

func (b *Instance) Create(instance *osbapi.Instance) error {
	instanceResource := &v1alpha1.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.ServiceInstanceSpec{
			Name:       instance.Name,
			PlanID:     instance.PlanID,
			ServiceID:  instance.ServiceID,
			BrokerName: instance.BrokerName,
		},
	}

	return b.KubeClient.Create(context.TODO(), instanceResource)
}
