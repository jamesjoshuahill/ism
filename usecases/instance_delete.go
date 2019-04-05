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

import "fmt"

type InstanceDeleteUsecase struct {
	InstanceDeleter InstanceDeleter
	InstanceFetcher InstanceFetcher
	BindingFetcher  BindingFetcher
}

func (u *InstanceDeleteUsecase) Delete(name string) error {
	instance, err := u.InstanceFetcher.GetInstanceByName(name)
	if err != nil {
		return err
	}

	bindings, err := u.BindingFetcher.GetBindingsForInstance(instance.ID)
	if err != nil {
		return err
	}

	if len(bindings) > 0 {
		return grammaticallyCorrectError(name, len(bindings))
	}

	return u.InstanceDeleter.Delete(name)
}

func grammaticallyCorrectError(name string, numBindings int) error {
	var pluralEnding string

	if numBindings > 1 {
		pluralEnding = "s"
	}

	return fmt.Errorf("Instance '%s' cannot be deleted as it has %d binding%s", name, numBindings, pluralEnding)
}
