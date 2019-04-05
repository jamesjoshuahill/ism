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

//go:generate counterfeiter . BindingCreator
//go:generate counterfeiter . BindingFetcher
//go:generate counterfeiter . BrokerFetcher
//go:generate counterfeiter . InstanceFetcher
//go:generate counterfeiter . InstanceCreator
//go:generate counterfeiter . InstanceDeleter
//go:generate counterfeiter . PlanFetcher
//go:generate counterfeiter . ServiceFetcher

type BindingCreator interface {
	Create(*osbapi.Binding) error
}

type BindingFetcher interface {
	GetBindingByName(name string) (*osbapi.Binding, error)
	GetBindingsForInstance(instanceID string) ([]*osbapi.Binding, error)
	GetBindings() ([]*osbapi.Binding, error)
}

type BrokerFetcher interface {
	GetBrokers() ([]*osbapi.Broker, error)
	GetBrokerByName(name string) (*osbapi.Broker, error)
}

type InstanceFetcher interface {
	GetInstanceByID(id string) (*osbapi.Instance, error)
	GetInstanceByName(name string) (*osbapi.Instance, error)
	GetInstances() ([]*osbapi.Instance, error)
}

type InstanceCreator interface {
	Create(*osbapi.Instance) error
}

type InstanceDeleter interface {
	Delete(name string) error
}

type PlanFetcher interface {
	GetPlanByID(id string) (*osbapi.Plan, error)
	GetPlans(serviceID string) ([]*osbapi.Plan, error)
}

type ServiceFetcher interface {
	GetServiceByID(id string) (*osbapi.Service, error)
	GetServices(brokerName string) ([]*osbapi.Service, error)
}
