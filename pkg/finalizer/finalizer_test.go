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

package finalizer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/pivotal-cf/ism/pkg/finalizer"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Finalizer", func() {
	DescribeTable("AddFinalizer",
		func(existingFinalizers []string, newFinalizer string, expectedFinalizers []string) {
			obj := &v1.ObjectMeta{}
			obj.SetFinalizers(existingFinalizers)

			AddFinalizer(obj, newFinalizer)
			Expect(obj.GetFinalizers()).To(Equal(expectedFinalizers))
		},
		Entry("no existing finalizers", []string{}, "finalizer.1", []string{"finalizer.1"}),
		Entry("one existing finalizer", []string{"finalizer.2"}, "finalizer.1", []string{"finalizer.2", "finalizer.1"}),
		Entry("some existing finalizers", []string{"finalizer.3", "finalizer.2"}, "finalizer.1", []string{"finalizer.3", "finalizer.2", "finalizer.1"}),
		Entry("finalizer exists", []string{"finalizer.1", "finalizer.2"}, "finalizer.1", []string{"finalizer.1", "finalizer.2"}),
	)

	DescribeTable("RemoveFinalizer",
		func(existingFinalizers []string, deleteFinalizer string, expectedFinalizers []string) {
			obj := &v1.ObjectMeta{}
			obj.SetFinalizers(existingFinalizers)

			RemoveFinalizer(obj, deleteFinalizer)
			Expect(obj.GetFinalizers()).To(Equal(expectedFinalizers))
		},
		Entry("no existing finalizers", []string{}, "finalizer.1", []string{}),
		Entry("one existing finalizer", []string{"finalizer.1"}, "finalizer.1", []string{}),
		Entry("some existing finalizers", []string{"finalizer.2", "finalizer.1"}, "finalizer.1", []string{"finalizer.2"}),
		Entry("no finalizer found", []string{"finalizer.3", "finalizer.2"}, "finalizer.1", []string{"finalizer.3", "finalizer.2"}),
	)
})
