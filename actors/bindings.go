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

package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . BindingRepository

type BindingRepository interface {
	FindAll() ([]*osbapi.Binding, error)
	Create(*osbapi.Binding) error
}

type BindingsActor struct {
	Repository BindingRepository
}

func (a *BindingsActor) Create(binding *osbapi.Binding) error {
	return a.Repository.Create(binding)
}

func (a *BindingsActor) GetBindings() ([]*osbapi.Binding, error) {
	return a.Repository.FindAll()
}
