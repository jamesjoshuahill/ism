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

//TODO: Where should this be defined?

//go:generate counterfeiter . UI

type UI interface {
	DisplayText(text string, data ...map[string]interface{})
	DisplayTable(table [][]string)
}

type RootCommand struct {
	BrokerCommand   BrokerCommand   `command:"broker" long-description:"The broker command group lets you register and list service brokers from the marketplace"`
	ServiceCommand  ServiceCommand  `command:"service" long-description:"The service command group lets you list the available services in the marketplace"`
	InstanceCommand InstanceCommand `command:"instance" long-description:"The instance command group lets you create and list service instances"`
	BindingCommand  BindingCommand  `command:"binding" long-description:"The binding command group lets you create service bindings"`
}
