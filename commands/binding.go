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

import "github.com/pivotal-cf/ism/usecases"

//go:generate counterfeiter . BindingCreateUsecase

type BindingCreateUsecase interface {
	Create(name, instanceName string) error
}

//go:generate counterfeiter . BindingListUsecase

type BindingListUsecase interface {
	GetBindings() ([]*usecases.Binding, error)
}

//go:generate counterfeiter . BindingGetUsecase

type BindingGetUsecase interface {
	GetBindingDetailsByName(name string) (*usecases.BindingDetails, error)
}

//go:generate counterfeiter . BindingDeleter

type BindingDeleter interface {
	Delete(name string) error
}

type BindingCommand struct {
	BindingCreateCommand BindingCreateCommand `command:"create" long-description:"Create a service binding"`
	BindingListCommand   BindingListCommand   `command:"list" long-description:"List the service bindings"`
	BindingGetCommand    BindingGetCommand    `command:"get" long-description:"Get a service binding"`
	BindingDeleteCommand BindingDeleteCommand `command:"delete" long-description:"Delete a service binding"`
}

type BindingCreateCommand struct {
	Name     string `long:"name" description:"Name of the service binding" required:"true"`
	Instance string `long:"instance" description:"Name of the service instance" required:"true"`

	UI                   UI
	BindingCreateUsecase BindingCreateUsecase
}

type BindingListCommand struct {
	UI                 UI
	BindingListUsecase BindingListUsecase
}

type BindingGetCommand struct {
	Name string `long:"name" description:"Name of the service binding" required:"true"`

	UI                UI
	BindingGetUsecase BindingGetUsecase
}

type BindingDeleteCommand struct {
	Name string `long:"name" description:"Name of the service binding" required:"true"`

	UI             UI
	BindingDeleter BindingDeleter
}

func (cmd *BindingCreateCommand) Execute([]string) error {
	if err := cmd.BindingCreateUsecase.Create(cmd.Name, cmd.Instance); err != nil {
		return err
	}

	cmd.UI.DisplayText("Binding '{{.BindingName}}' is being created.", map[string]interface{}{"BindingName": cmd.Name})
	return nil
}

func (cmd *BindingListCommand) Execute([]string) error {
	bindings, err := cmd.BindingListUsecase.GetBindings()
	if err != nil {
		return err
	}

	if len(bindings) == 0 {
		cmd.UI.DisplayText("No bindings found.")
		return nil
	}

	bindingsTable := buildBindingTableData(bindings)
	cmd.UI.DisplayTable(bindingsTable)
	return nil
}

func (cmd *BindingGetCommand) Execute([]string) error {
	binding, err := cmd.BindingGetUsecase.GetBindingDetailsByName(cmd.Name)
	if err != nil {
		return err
	}

	return cmd.UI.DisplayYAML(binding)
}

func (cmd *BindingDeleteCommand) Execute([]string) error {
	if err := cmd.BindingDeleter.Delete(cmd.Name); err != nil {
		return err
	}

	cmd.UI.DisplayText("Binding '{{.BindingName}}' is being deleted.", map[string]interface{}{"BindingName": cmd.Name})
	return nil
}

func buildBindingTableData(bindings []*usecases.Binding) [][]string {
	headers := []string{"NAME", "INSTANCE", "STATUS", "CREATED AT"}
	data := [][]string{headers}

	for _, b := range bindings {
		row := []string{b.Name, b.InstanceName, b.Status, b.CreatedAt}
		data = append(data, row)
	}
	return data
}
