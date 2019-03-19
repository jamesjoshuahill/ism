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

//go:generate counterfeiter . BindingCreateUsecase

type BindingCreateUsecase interface {
	Create(name, instanceName string) error
}

type BindingCommand struct {
	BindingCreateCommand BindingCreateCommand `command:"create" long-description:"Create a service binding"`
}

type BindingCreateCommand struct {
	Name     string `long:"name" description:"Name of the service binding" required:"true"`
	Instance string `long:"instance" description:"Name of the service instance" required:"true"`

	UI                   UI
	BindingCreateUsecase BindingCreateUsecase
}

func (cmd *BindingCreateCommand) Execute([]string) error {
	if err := cmd.BindingCreateUsecase.Create(cmd.Name, cmd.Instance); err != nil {
		return err
	}

	cmd.UI.DisplayText("Binding '{{.BindingName}}' is being created.", map[string]interface{}{"BindingName": cmd.Name})
	return nil
}
