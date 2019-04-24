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
	"github.com/pivotal-cf/ism/usecases"
)

//go:generate counterfeiter . InstanceCreateUsecase

type InstanceCreateUsecase interface {
	Create(name, planName, serviceName, brokerName string) error
}

//go:generate counterfeiter . InstanceGetUsecase

type InstanceGetUsecase interface {
	GetInstanceDetailsByName(name string) (*usecases.InstanceDetails, error)
}

//go:generate counterfeiter . InstanceListUsecase

type InstanceListUsecase interface {
	GetInstances() ([]*usecases.Instance, error)
}

//go:generate counterfeiter . InstanceDeleteUsecase

type InstanceDeleteUsecase interface {
	Delete(name string) error
}

type InstanceCommand struct {
	InstanceCreateCommand InstanceCreateCommand `command:"create" long-description:"Create a service instance"`
	InstanceListCommand   InstanceListCommand   `command:"list" long-description:"List the service instances"`
	InstanceDeleteCommand InstanceDeleteCommand `command:"delete" long-description:"Delete a service instance"`
	InstanceGetCommand    InstanceGetCommand    `command:"get" long-description:"Get a service instance"`
}

type InstanceCreateCommand struct {
	Name    string `long:"name" description:"Name of the service instance" required:"true"`
	Service string `long:"service" description:"Name of the service" required:"true"`
	Plan    string `long:"plan" description:"Name of the plan" required:"true"`
	Broker  string `long:"broker" description:"Name of the broker" required:"true"`

	UI                    UI
	InstanceCreateUsecase InstanceCreateUsecase
}

type InstanceGetCommand struct {
	Name string `long:"name" description:"Name of the service instance" required:"true"`

	UI                 UI
	InstanceGetUsecase InstanceGetUsecase
}

type InstanceListCommand struct {
	UI                  UI
	InstanceListUsecase InstanceListUsecase
}

type InstanceDeleteCommand struct {
	Name string `long:"name" description:"Name of the service instance" required:"true"`

	UI                    UI
	InstanceDeleteUsecase InstanceDeleteUsecase
}

func (cmd *InstanceCreateCommand) Execute([]string) error {
	if err := cmd.InstanceCreateUsecase.Create(cmd.Name, cmd.Plan, cmd.Service, cmd.Broker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Instance '{{.InstanceName}}' is being created.", map[string]interface{}{"InstanceName": cmd.Name})

	return nil
}

func (cmd *InstanceGetCommand) Execute([]string) error {
	instance, err := cmd.InstanceGetUsecase.GetInstanceDetailsByName(cmd.Name)
	if err != nil {
		return err
	}
	return cmd.UI.DisplayYAML(instance)
}

func (cmd *InstanceListCommand) Execute([]string) error {
	instances, err := cmd.InstanceListUsecase.GetInstances()
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		cmd.UI.DisplayText("No instances found.")
		return nil
	}

	instancesTable := buildInstanceTableData(instances)
	cmd.UI.DisplayTable(instancesTable)
	return nil
}

func (cmd *InstanceDeleteCommand) Execute([]string) error {
	if err := cmd.InstanceDeleteUsecase.Delete(cmd.Name); err != nil {
		return err
	}

	cmd.UI.DisplayText("Instance '{{.InstanceName}}' is being deleted.", map[string]interface{}{"InstanceName": cmd.Name})
	return nil
}

func buildInstanceTableData(instances []*usecases.Instance) [][]string {
	headers := []string{"NAME", "SERVICE", "PLAN", "BROKER", "STATUS", "CREATED AT"}
	data := [][]string{headers}

	for _, i := range instances {
		row := []string{i.Name, i.ServiceName, i.PlanName, i.BrokerName, i.Status, i.CreatedAt}
		data = append(data, row)
	}

	return data
}
