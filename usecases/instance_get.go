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

type InstanceDetails struct {
	Name        string `json:"name"`
	PlanName    string `json:"plan"`
	ServiceName string `json:"service"`
	BrokerName  string `json:"broker"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

type InstanceGetUsecase struct {
	InstanceFetcher InstanceFetcher
	ServiceFetcher  ServiceFetcher
	PlanFetcher     PlanFetcher
	BrokerFetcher   BrokerFetcher
}

func (usecase *InstanceGetUsecase) GetInstanceDetailsByName(name string) (*InstanceDetails, error) {
	instance, err := usecase.InstanceFetcher.GetInstanceByName(name)
	if err != nil {
		return &InstanceDetails{}, err
	}

	service, err := usecase.ServiceFetcher.GetServiceByID(instance.ServiceID)
	if err != nil {
		return &InstanceDetails{}, err
	}

	plan, err := usecase.PlanFetcher.GetPlanByID(instance.PlanID)
	if err != nil {
		return &InstanceDetails{}, err
	}

	broker, err := usecase.BrokerFetcher.GetBrokerByName(instance.BrokerName)
	if err != nil {
		return &InstanceDetails{}, err
	}

	return &InstanceDetails{
		Name:        instance.Name,
		ServiceName: service.Name,
		PlanName:    plan.Name,
		BrokerName:  broker.Name,
		Status:      instance.Status,
		CreatedAt:   instance.CreatedAt,
	}, nil
}
