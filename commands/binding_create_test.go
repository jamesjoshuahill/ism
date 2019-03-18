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

package commands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
)

var _ = Describe("Binding create command", func() {

	var (
		fakeBindingCreateUsecase *commandsfakes.FakeBindingCreateUsecase
		fakeUI                   *commandsfakes.FakeUI

		createCommand BindingCreateCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeBindingCreateUsecase = &commandsfakes.FakeBindingCreateUsecase{}
		fakeUI = &commandsfakes.FakeUI{}

		createCommand = BindingCreateCommand{
			Name:                 "my-binding-1",
			InstanceName:         "my-instance-1",
			BindingCreateUsecase: fakeBindingCreateUsecase,
			UI:                   fakeUI,
		}
	})

	JustBeforeEach(func() {
		executeErr = createCommand.Execute(nil)
	})

	When("creating a binding succeeds", func() {
		BeforeEach(func() {
			fakeBindingCreateUsecase.CreateReturns(nil)
		})

		It("calls to create the binding", func() {
			name, instanceName := fakeBindingCreateUsecase.CreateArgsForCall(0)

			Expect(name).To(Equal("my-binding-1"))
			Expect(instanceName).To(Equal("my-instance-1"))
		})

		It("displays that the service binding is being created", func() {
			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Binding '{{.BindingName}}' is being created."))
			Expect(data[0]).To(HaveKeyWithValue("BindingName", "my-binding-1"))
		})
	})

	When("creating a binding errors", func() {
		BeforeEach(func() {
			fakeBindingCreateUsecase.CreateReturns(errors.New("error-creating-binding"))
		})

		It("returns the error", func() {
			Expect(executeErr).To(MatchError("error-creating-binding"))
		})
	})
})
