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

package usecases

import (
	"github.com/pivotal-cf/ism/osbapi"
)

type Instance struct {
	Name        string
	ServiceName string
	PlanName    string
	BrokerName  string
	Status      string
	CreatedAt   string
}

//go:generate counterfeiter . InstancesFetcher

type InstancesFetcher interface {
	GetInstances() ([]*osbapi.Instance, error)
}

//go:generate counterfeiter . ServiceFetcher

type ServiceFetcher interface {
	GetService(serviceID string) (*osbapi.Service, error)
}

//go:generate counterfeiter . PlanFetcher

type PlanFetcher interface {
	GetPlan(planID string) (*osbapi.Plan, error)
}

type InstanceListUsecase struct {
	InstancesFetcher InstancesFetcher
	ServiceFetcher   ServiceFetcher
	PlanFetcher      PlanFetcher
}

func (i *InstanceListUsecase) GetInstances() ([]*Instance, error) {
	osbapiInstances, err := i.InstancesFetcher.GetInstances()
	if err != nil {
		return []*Instance{}, err
	}

	var instances []*Instance
	for _, instance := range osbapiInstances {
		service, err := i.ServiceFetcher.GetService(instance.ServiceID)
		if err != nil {
			return []*Instance{}, err
		}

		plan, err := i.PlanFetcher.GetPlan(instance.PlanID)
		if err != nil {
			return []*Instance{}, err
		}

		instances = append(instances, &Instance{
			Name:        instance.Name,
			ServiceName: service.Name,
			PlanName:    plan.Name,
			BrokerName:  instance.BrokerName,
			Status:      instance.Status,
			CreatedAt:   instance.CreatedAt,
		})
	}

	return instances, nil
}
