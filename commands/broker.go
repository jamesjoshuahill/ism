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
	"github.com/pivotal-cf/ism/osbapi"
)

//go:generate counterfeiter . BrokerRegistrar

type BrokerRegistrar interface {
	Register(*osbapi.Broker) error
}

//go:generate counterfeiter . BrokersFetcher

type BrokersFetcher interface {
	GetBrokers() ([]*osbapi.Broker, error)
}

//go:generate counterfeiter . BrokerDeleteUsecase

type BrokerDeleteUsecase interface {
	Delete(name string) error
}

type BrokerCommand struct {
	BrokerRegisterCommand BrokerRegisterCommand `command:"register" long-description:"Register a service broker into the marketplace"`
	BrokerListCommand     BrokerListCommand     `command:"list" long-description:"List the service brokers in the marketplace"`
	BrokerDeleteCommand   BrokerDeleteCommand   `command:"delete" long-description:"Delete a service broker from the marketplace"`
}

type BrokerListCommand struct {
	UI             UI
	BrokersFetcher BrokersFetcher
}

type BrokerRegisterCommand struct {
	Name     string `long:"name" description:"Name of the service broker" required:"true"`
	URL      string `long:"url" description:"URL of the service broker" required:"true"`
	Username string `long:"username" description:"Username of the service broker" required:"true"`
	Password string `long:"password" description:"Password of the service broker" required:"true"`

	UI              UI
	BrokerRegistrar BrokerRegistrar
}

type BrokerDeleteCommand struct {
	Name string `long:"name" description:"Name of the service broker" required:"true"`

	UI                  UI
	BrokerDeleteUsecase BrokerDeleteUsecase
}

func (cmd *BrokerRegisterCommand) Execute([]string) error {
	//TODO: This is the only command that uses the osbapi types, should this just pass params instead?
	newBroker := &osbapi.Broker{
		Name:     cmd.Name,
		URL:      cmd.URL,
		Username: cmd.Username,
		Password: cmd.Password,
	}

	if err := cmd.BrokerRegistrar.Register(newBroker); err != nil {
		return err
	}

	cmd.UI.DisplayText("Broker '{{.BrokerName}}' registered.", map[string]interface{}{"BrokerName": cmd.Name})

	return nil
}

func (cmd *BrokerListCommand) Execute([]string) error {
	brokers, err := cmd.BrokersFetcher.GetBrokers()
	if err != nil {
		return err
	}

	if len(brokers) == 0 {
		cmd.UI.DisplayText("No brokers found.")
		return nil
	}

	brokersTable := buildBrokerTableData(brokers)
	cmd.UI.DisplayTable(brokersTable)
	return nil
}

func (cmd *BrokerDeleteCommand) Execute([]string) error {
	if err := cmd.BrokerDeleteUsecase.Delete(cmd.Name); err != nil {
		return err
	}

	cmd.UI.DisplayText("Broker '{{.BrokerName}}' is being deleted.", map[string]interface{}{"BrokerName": cmd.Name})
	return nil
}

func buildBrokerTableData(brokers []*osbapi.Broker) [][]string {
	headers := []string{"NAME", "URL", "CREATED AT"}
	data := [][]string{headers}

	for _, b := range brokers {
		row := []string{b.Name, b.URL, b.CreatedAt}
		data = append(data, row)
	}
	return data
}
