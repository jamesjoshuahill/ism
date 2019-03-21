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

type Binding struct {
	Name         string
	InstanceName string
	Status       string
	CreatedAt    string
}

type BindingListUsecase struct {
	BindingFetcher  BindingFetcher
	InstanceFetcher InstanceFetcher
}

func (b *BindingListUsecase) GetBindings() ([]*Binding, error) {
	osbapiBindings, err := b.BindingFetcher.GetBindings()
	if err != nil {
		return []*Binding{}, err
	}

	var bindings []*Binding
	for _, binding := range osbapiBindings {
		instance, err := b.InstanceFetcher.GetInstanceByID(binding.InstanceID)
		if err != nil {
			return []*Binding{}, err
		}

		bindings = append(bindings, &Binding{
			Name:         binding.Name,
			InstanceName: instance.Name,
			Status:       binding.Status,
			CreatedAt:    binding.CreatedAt,
		})
	}

	return bindings, nil
}
