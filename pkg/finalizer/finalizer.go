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

package finalizer

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func AddFinalizer(obj v1.Object, finalizer string) {
	finalizers := obj.GetFinalizers()
	for _, f := range finalizers {
		if f == finalizer {
			return
		}
	}

	obj.SetFinalizers(append(finalizers, finalizer))
	return
}

func RemoveFinalizer(obj v1.Object, finalizer string) {
	finalizers := obj.GetFinalizers()

	newFinalizers := []string{}
	for _, f := range finalizers {
		if f != finalizer {
			newFinalizers = append(newFinalizers, f)
		}
	}

	obj.SetFinalizers(newFinalizers)
	return
}
