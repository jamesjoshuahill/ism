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

type BindingDetails struct {
	Name         string                 `json:"name"`
	InstanceName string                 `json:"instance"`
	PlanName     string                 `json:"plan"`
	ServiceName  string                 `json:"service"`
	BrokerName   string                 `json:"broker"`
	Status       string                 `json:"status"`
	CreatedAt    string                 `json:"createdAt"`
	Credentials  map[string]interface{} `json:"credentials"`
}

type BindingGetUsecase struct {
	BindingFetcher  BindingFetcher
	InstanceFetcher InstanceFetcher
	ServiceFetcher  ServiceFetcher
	PlanFetcher     PlanFetcher
	BrokerFetcher   BrokerFetcher
}

func (b *BindingGetUsecase) GetBindingDetailsByName(name string) (*BindingDetails, error) {
	osbapiBinding, err := b.BindingFetcher.GetBindingByName(name)
	if err != nil {
		return &BindingDetails{}, err
	}

	osbapiInstance, err := b.InstanceFetcher.GetInstanceByID(osbapiBinding.InstanceID)
	if err != nil {
		return &BindingDetails{}, err
	}

	osbapiService, err := b.ServiceFetcher.GetServiceByID(osbapiBinding.ServiceID)
	if err != nil {
		return &BindingDetails{}, err
	}

	osbapiPlan, err := b.PlanFetcher.GetPlanByID(osbapiBinding.PlanID)
	if err != nil {
		return &BindingDetails{}, err
	}

	osbapiBroker, err := b.BrokerFetcher.GetBrokerByName(osbapiBinding.BrokerName)
	if err != nil {
		return &BindingDetails{}, err
	}

	return &BindingDetails{
		Name:         osbapiBinding.Name,
		InstanceName: osbapiInstance.Name,
		ServiceName:  osbapiService.Name,
		PlanName:     osbapiPlan.Name,
		BrokerName:   osbapiBroker.Name,
		Status:       osbapiBinding.Status,
		CreatedAt:    osbapiBinding.CreatedAt,
		Credentials:  osbapiBinding.Credentials,
	}, err
}
