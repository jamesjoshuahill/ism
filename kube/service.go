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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type Service struct {
	KubeClient client.Client
}

func (s *Service) Find(serviceID string) (*osbapi.Service, error) {
	service := &v1alpha1.BrokerService{}
	err := s.KubeClient.Get(context.TODO(), types.NamespacedName{Name: serviceID, Namespace: "default"}, service)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &osbapi.Service{
		ID:          service.ObjectMeta.Name,
		Name:        service.Spec.Name,
		Description: service.Spec.Description,
		BrokerName:  service.Spec.BrokerID,
	}, nil
}

func (s *Service) FindByBroker(brokerName string) ([]*osbapi.Service, error) {
	list := &v1alpha1.BrokerServiceList{}
	if err := s.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Service{}, err
	}

	services := []*osbapi.Service{}
	for _, s := range list.Items {
		// TODO: This code will be refactored so filtering happens in the API - for now
		// we are assuming there will never be multiple owner references. See #164327846
		if len(s.ObjectMeta.OwnerReferences) == 0 {
			break
		}

		if string(s.ObjectMeta.OwnerReferences[0].Name) == brokerName {
			services = append(services, &osbapi.Service{
				ID:          s.ObjectMeta.Name,
				Name:        s.Spec.Name,
				Description: s.Spec.Description,
				BrokerName:  brokerName,
			})
		}
	}

	return services, nil
}
