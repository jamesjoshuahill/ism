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
	"github.com/pivotal-cf/ism/usecases"
)

var _ = Describe("Binding command", func() {
	Describe("Create sub command", func() {
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
				Instance:             "my-instance-1",
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

	Describe("Get sub command", func() {
		var (
			fakeBindingGetUsecase *commandsfakes.FakeBindingGetUsecase
			fakeUI                *commandsfakes.FakeUI

			getCommand BindingGetCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeBindingGetUsecase = &commandsfakes.FakeBindingGetUsecase{}
			fakeUI = &commandsfakes.FakeUI{}

			getCommand = BindingGetCommand{
				Name:              "my-binding",
				BindingGetUsecase: fakeBindingGetUsecase,
				UI:                fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = getCommand.Execute(nil)
		})

		It("calls to get the binding", func() {
			bindingName := fakeBindingGetUsecase.GetBindingDetailsByNameArgsForCall(0)
			Expect(bindingName).To(Equal("my-binding"))
		})

		When("the binding exists", func() {
			BeforeEach(func() {
				fakeBindingGetUsecase.GetBindingDetailsByNameReturns(&usecases.BindingDetails{Name: "my-binding"}, nil)
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("passes the binding details for display", func() {
				Expect(fakeUI.DisplayYAMLCallCount()).To(Equal(1))

				bindingDetailsArg := fakeUI.DisplayYAMLArgsForCall(0)

				Expect(bindingDetailsArg).To(Equal(&usecases.BindingDetails{Name: "my-binding"}))
			})
		})

		When("get binding by name fails", func() {
			BeforeEach(func() {
				fakeBindingGetUsecase.GetBindingDetailsByNameReturns(&usecases.BindingDetails{}, errors.New("error-binding-not-found"))
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError("error-binding-not-found"))
			})
		})

		When("displaying the binding errors", func() {
			BeforeEach(func() {
				fakeUI.DisplayYAMLReturns(errors.New("error-displaying-yaml"))
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError("error-displaying-yaml"))
			})
		})
	})

	Describe("List sub command", func() {
		var (
			fakeBindingListUsecase *commandsfakes.FakeBindingListUsecase
			fakeUI                 *commandsfakes.FakeUI

			listCommand BindingListCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeBindingListUsecase = &commandsfakes.FakeBindingListUsecase{}
			fakeUI = &commandsfakes.FakeUI{}

			listCommand = BindingListCommand{
				BindingListUsecase: fakeBindingListUsecase,
				UI:                 fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = listCommand.Execute(nil)
		})

		When("there are no bindings", func() {
			BeforeEach(func() {
				fakeBindingListUsecase.GetBindingsReturns([]*usecases.Binding{}, nil)
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("displays that no bindings were found", func() {
				Expect(fakeUI.DisplayTextCallCount()).NotTo(BeZero())
				text, _ := fakeUI.DisplayTextArgsForCall(0)

				Expect(text).To(Equal("No bindings found."))
			})
		})

		When("there are 1 or more bindings", func() {
			BeforeEach(func() {
				fakeBindingListUsecase.GetBindingsReturns([]*usecases.Binding{
					{
						Name:         "my-binding-1",
						InstanceName: "my-instance-1",
						CreatedAt:    "some-time-1",
						Status:       "created",
					},
					{
						Name:         "my-binding-2",
						InstanceName: "my-instance-2",
						CreatedAt:    "some-time-2",
						Status:       "creating",
					}}, nil)
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("displays the bindings in a table", func() {
				Expect(fakeUI.DisplayTableCallCount()).NotTo(BeZero())
				data := fakeUI.DisplayTableArgsForCall(0)

				Expect(data[0]).To(Equal([]string{"NAME", "INSTANCE", "STATUS", "CREATED AT"}))
				Expect(data[1]).To(Equal([]string{"my-binding-1", "my-instance-1", "created", "some-time-1"}))
				Expect(data[2]).To(Equal([]string{"my-binding-2", "my-instance-2", "creating", "some-time-2"}))
			})
		})

		When("getting services returns an error", func() {
			BeforeEach(func() {
				fakeBindingListUsecase.GetBindingsReturns([]*usecases.Binding{}, errors.New("error-getting-bindings"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-getting-bindings"))
			})
		})
	})
})
