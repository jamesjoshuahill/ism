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

import "github.com/pivotal-cf/ism/osbapi"

type Service struct {
	Name        string
	Description string
	PlanNames   []string
	BrokerName  string
}

type ServiceListUsecase struct {
	BrokerFetcher  BrokerFetcher
	ServiceFetcher ServiceFetcher
	PlanFetcher    PlanFetcher
}

func (u *ServiceListUsecase) GetServices() ([]*Service, error) {
	brokers, err := u.BrokerFetcher.GetBrokers()
	if err != nil {
		return []*Service{}, err
	}

	var services []*Service
	for _, b := range brokers {
		osbapiServices, err := u.ServiceFetcher.GetServices(b.Name)
		if err != nil {
			return []*Service{}, err
		}

		for _, osbapiService := range osbapiServices {
			plans, err := u.PlanFetcher.GetPlans(osbapiService.ID)
			if err != nil {
				return []*Service{}, err
			}

			service := &Service{
				Name:        osbapiService.Name,
				Description: osbapiService.Description,
				PlanNames:   plansToNames(plans),
				BrokerName:  b.Name,
			}
			services = append(services, service)
		}
	}

	return services, nil
}

func plansToNames(plans []*osbapi.Plan) []string {
	names := []string{}
	for _, p := range plans {
		names = append(names, p.Name)
	}

	return names
}
