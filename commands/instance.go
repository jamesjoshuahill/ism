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

//go:generate counterfeiter . InstanceCreateUsecase

type InstanceCreateUsecase interface {
	Create(name, planName, serviceName, brokerName string) error
}

type InstanceCommand struct {
	InstanceCreateCommand InstanceCreateCommand `command:"create" long-description:"Create a service instance"`
}

type InstanceCreateCommand struct {
	Name    string `long:"name" description:"Name of the service instance" required:"true"`
	Service string `long:"service" description:"Name of the service" required:"true"`
	Plan    string `long:"plan" description:"Name of the plan" required:"true"`
	Broker  string `long:"broker" description:"Name of the broker" required:"true"`

	UI                    UI
	InstanceCreateUsecase InstanceCreateUsecase
}

func (cmd *InstanceCreateCommand) Execute([]string) error {
	if err := cmd.InstanceCreateUsecase.Create(cmd.Name, cmd.Plan, cmd.Service, cmd.Broker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Instance '{{.InstanceName}}' created.", map[string]interface{}{"InstanceName": cmd.Name})

	return nil
}
