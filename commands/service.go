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

package commands

import (
	"strings"

	"github.com/pivotal-cf/ism/usecases"
)

//go:generate counterfeiter . ServiceListUsecase

type ServiceListUsecase interface {
	GetServices() ([]*usecases.Service, error)
}

type ServiceCommand struct {
	ServiceListCommand ServiceListCommand `command:"list" long-description:"List the services that are available in the marketplace."`
}

type ServiceListCommand struct {
	UI                 UI
	ServiceListUsecase ServiceListUsecase
}

func (cmd *ServiceListCommand) Execute([]string) error {
	services, err := cmd.ServiceListUsecase.GetServices()
	if err != nil {
		return err
	}

	if len(services) == 0 {
		cmd.UI.DisplayText("No services found.")
		return nil
	}

	servicesTable := buildServiceTableData(services)
	cmd.UI.DisplayTable(servicesTable)

	return nil
}

func buildServiceTableData(services []*usecases.Service) [][]string {
	headers := []string{"SERVICE", "PLANS", "BROKER", "DESCRIPTION"}
	data := [][]string{headers}

	for _, s := range services {
		row := []string{s.Name, strings.Join(s.PlanNames, ", "), s.BrokerName, s.Description}
		data = append(data, row)
	}

	return data
}
