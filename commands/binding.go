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

import "fmt"

type BindingCommand struct {
	BindingCreateCommand BindingCreateCommand `command:"create" long-description:"Create a service binding"`
}

type BindingCreateCommand struct {
	Name         string `long:"name" description:"Name of the service binding" required:"true"`
	InstanceName string `long:"instance-name" description:"Name of the service instance" required:"true"`
}

func (cmd *BindingCreateCommand) Execute([]string) error {
	fmt.Println("Binding 'my-binding' is being created.")
	return nil
}
