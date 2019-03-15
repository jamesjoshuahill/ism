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
	"fmt"

	"github.com/pivotal-cf/ism/osbapi"
)

//go:generate counterfeiter . InstanceCreator

type InstanceCreator interface {
	Create(*osbapi.Instance) error
}

type InstanceCreateUsecase struct {
	BrokersFetcher  BrokersFetcher
	ServicesFetcher ServicesFetcher
	PlansFetcher    PlansFetcher
	InstanceCreator InstanceCreator
}

func (u *InstanceCreateUsecase) Create(name, planName, serviceName, brokerName string) error {
	broker, err := u.getBroker(brokerName)
	if err != nil {
		return err
	}

	service, err := u.getService(broker.Name, serviceName)
	if err != nil {
		return err
	}

	plan, err := u.getPlan(service.ID, planName)
	if err != nil {
		return err
	}

	instance := &osbapi.Instance{
		Name:       name,
		PlanID:     plan.ID,
		ServiceID:  service.ID,
		BrokerName: broker.Name,
	}

	return u.InstanceCreator.Create(instance)
}

func (u *InstanceCreateUsecase) getBroker(brokerName string) (*osbapi.Broker, error) {
	brokers, err := u.BrokersFetcher.GetBrokers()
	if err != nil {
		return &osbapi.Broker{}, err
	}

	// TODO: This code will be refactored so filtering happens in the API. See #164327846
	for _, broker := range brokers {
		if broker.Name == brokerName {
			return broker, nil
		}
	}

	return &osbapi.Broker{}, fmt.Errorf("Broker '%s' does not exist", brokerName)
}

func (u *InstanceCreateUsecase) getService(brokerName, serviceName string) (*osbapi.Service, error) {
	services, err := u.ServicesFetcher.GetServices(brokerName)
	if err != nil {
		return &osbapi.Service{}, err
	}

	// TODO: This code will be refactored so filtering happens in the API. See #164327846
	for _, service := range services {
		if service.Name == serviceName {
			return service, nil
		}
	}

	return &osbapi.Service{}, fmt.Errorf("Service '%s' does not exist", serviceName)
}

func (u *InstanceCreateUsecase) getPlan(serviceID, planName string) (*osbapi.Plan, error) {
	plans, err := u.PlansFetcher.GetPlans(serviceID)
	if err != nil {
		return &osbapi.Plan{}, err
	}

	// TODO: This code will be refactored so filtering happens in the API. See #164327846
	for _, plan := range plans {
		if plan.Name == planName {
			return plan, nil
		}
	}

	return &osbapi.Plan{}, fmt.Errorf("Plan '%s' does not exist", planName)
}
