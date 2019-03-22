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

type Instance struct {
	Name        string
	ServiceName string
	PlanName    string
	BrokerName  string
	Status      string
	CreatedAt   string
}

type InstanceListUsecase struct {
	InstanceFetcher InstanceFetcher
	ServiceFetcher  ServiceFetcher
	PlanFetcher     PlanFetcher
}

func (i *InstanceListUsecase) GetInstances() ([]*Instance, error) {
	osbapiInstances, err := i.InstanceFetcher.GetInstances()
	if err != nil {
		return []*Instance{}, err
	}

	var instances []*Instance
	for _, instance := range osbapiInstances {
		service, err := i.ServiceFetcher.GetServiceByID(instance.ServiceID)
		if err != nil {
			return []*Instance{}, err
		}

		plan, err := i.PlanFetcher.GetPlanByID(instance.PlanID)
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
