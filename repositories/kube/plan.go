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
	"errors"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

var errPlanNotFound = errors.New("plan not found")

type Plan struct {
	KubeClient client.Client
}

func (p *Plan) GetPlanByID(planID string) (*osbapi.Plan, error) {
	plan := &v1alpha1.BrokerServicePlan{}
	err := p.KubeClient.Get(context.TODO(), types.NamespacedName{Name: planID, Namespace: "default"}, plan)

	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, errPlanNotFound
		}
		return nil, err
	}

	return &osbapi.Plan{
		ID:        plan.ObjectMeta.Name,
		Name:      plan.Spec.Name,
		ServiceID: plan.Spec.ServiceID,
	}, nil
}

func (p *Plan) GetPlans(serviceID string) ([]*osbapi.Plan, error) {
	list := &v1alpha1.BrokerServicePlanList{}
	if err := p.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Plan{}, err
	}

	plans := []*osbapi.Plan{}
	for _, p := range list.Items {
		// TODO: This code will be refactored so filtering happens in the API - for now
		// we are assuming there will never be multiple owner references. See #164327846
		if len(p.ObjectMeta.OwnerReferences) == 0 {
			break
		}

		if p.ObjectMeta.OwnerReferences[0].Name == serviceID {
			plans = append(plans, &osbapi.Plan{
				ID:        p.ObjectMeta.Name,
				Name:      p.Spec.Name,
				ServiceID: serviceID,
			})
		}
	}

	return plans, nil
}
