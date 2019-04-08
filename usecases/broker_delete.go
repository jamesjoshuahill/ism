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

import (
	"errors"
	"fmt"
)

type BrokerDeleteUsecase struct {
	BrokerDeleter   BrokerDeleter
	InstanceFetcher InstanceFetcher
}

func (u *BrokerDeleteUsecase) Delete(name string) error {
	instances, err := u.InstanceFetcher.GetInstancesForBroker(name)
	if err != nil {
		return err
	}

	if len(instances) > 0 {
		errorMessage := fmt.Sprintf("Broker '%s' cannot be deleted as it has %d instance", name, len(instances))
		return grammaticallyCorrectError(errorMessage, len(instances))
	}

	return u.BrokerDeleter.Delete(name)
}

func grammaticallyCorrectError(errorMessage string, num int) error {
	if num > 1 {
		errorMessage = errorMessage + "s"
	}

	return errors.New(errorMessage)
}
