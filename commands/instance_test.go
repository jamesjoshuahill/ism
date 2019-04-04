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

var _ = Describe("Instance Command", func() {
	Describe("Create sub command", func() {
		var (
			fakeInstanceCreateUsecase *commandsfakes.FakeInstanceCreateUsecase
			fakeUI                    *commandsfakes.FakeUI

			createCommand InstanceCreateCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeInstanceCreateUsecase = &commandsfakes.FakeInstanceCreateUsecase{}
			fakeUI = &commandsfakes.FakeUI{}

			createCommand = InstanceCreateCommand{
				Name:                  "instance-1",
				Plan:                  "plan-1",
				Service:               "service-1",
				Broker:                "broker-1",
				InstanceCreateUsecase: fakeInstanceCreateUsecase,
				UI:                    fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = createCommand.Execute(nil)
		})

		When("creating an instance succeeds", func() {
			BeforeEach(func() {
				fakeInstanceCreateUsecase.CreateReturns(nil)
			})

			It("calls to create the instance", func() {
				name, planName, serviceName, brokerName := fakeInstanceCreateUsecase.CreateArgsForCall(0)

				Expect(name).To(Equal("instance-1"))
				Expect(planName).To(Equal("plan-1"))
				Expect(serviceName).To(Equal("service-1"))
				Expect(brokerName).To(Equal("broker-1"))
			})

			It("displays that the service instance is being created", func() {
				text, data := fakeUI.DisplayTextArgsForCall(0)
				Expect(text).To(Equal("Instance '{{.InstanceName}}' is being created."))
				Expect(data[0]).To(HaveKeyWithValue("InstanceName", "instance-1"))
			})
		})

		When("creating an instance errors", func() {
			BeforeEach(func() {
				fakeInstanceCreateUsecase.CreateReturns(errors.New("error-creating-instance"))
			})

			It("returns the error", func() {
				Expect(executeErr).To(MatchError("error-creating-instance"))
			})
		})
	})

	Describe("List sub command", func() {
		var (
			fakeUI                  *commandsfakes.FakeUI
			fakeInstanceListUsecase *commandsfakes.FakeInstanceListUsecase

			listCommand InstanceListCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeInstanceListUsecase = &commandsfakes.FakeInstanceListUsecase{}
			fakeUI = &commandsfakes.FakeUI{}

			listCommand = InstanceListCommand{
				InstanceListUsecase: fakeInstanceListUsecase,
				UI:                  fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = listCommand.Execute(nil)
		})

		When("there are no instances", func() {
			BeforeEach(func() {
				fakeInstanceListUsecase.GetInstancesReturns([]*usecases.Instance{}, nil)
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("displays that no instances were found", func() {
				Expect(fakeUI.DisplayTextCallCount()).NotTo(BeZero())
				text, _ := fakeUI.DisplayTextArgsForCall(0)

				Expect(text).To(Equal("No instances found."))
			})
		})

		When("there are 1 or more instance", func() {
			BeforeEach(func() {
				fakeInstanceListUsecase.GetInstancesReturns([]*usecases.Instance{{
					Name:        "my-instance-1",
					ServiceName: "my-service-1",
					PlanName:    "my-plan-1",
					BrokerName:  "my-broker-1",
					Status:      "created",
					CreatedAt:   "2019-02-28T12:08:31Z",
				}, {
					Name:        "my-instance-2",
					ServiceName: "my-service-1",
					PlanName:    "my-plan-2",
					BrokerName:  "my-broker-1",
					Status:      "creating",
					CreatedAt:   "2019-02-28T12:08:31Z",
				}}, nil)
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("displays the instances in a table ", func() {
				Expect(fakeUI.DisplayTableCallCount()).NotTo(BeZero())
				data := fakeUI.DisplayTableArgsForCall(0)

				Expect(data[0]).To(Equal([]string{"NAME", "SERVICE", "PLAN", "BROKER", "STATUS", "CREATED AT"}))
				Expect(data[1]).To(Equal([]string{"my-instance-1", "my-service-1", "my-plan-1", "my-broker-1", "created", "2019-02-28T12:08:31Z"}))
				Expect(data[2]).To(Equal([]string{"my-instance-2", "my-service-1", "my-plan-2", "my-broker-1", "creating", "2019-02-28T12:08:31Z"}))
			})
		})

		When("getting instances returns an error", func() {
			BeforeEach(func() {
				fakeInstanceListUsecase.GetInstancesReturns([]*usecases.Instance{}, errors.New("error-getting-instances"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-getting-instances"))
			})
		})
	})

	Describe("Delete sub command", func() {
		var (
			fakeInstanceDeleter *commandsfakes.FakeInstanceDeleter
			fakeUI              *commandsfakes.FakeUI

			deleteCommand InstanceDeleteCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeInstanceDeleter = &commandsfakes.FakeInstanceDeleter{}
			fakeUI = &commandsfakes.FakeUI{}

			deleteCommand = InstanceDeleteCommand{
				Name:            "my-instance",
				InstanceDeleter: fakeInstanceDeleter,
				UI:              fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = deleteCommand.Execute(nil)
		})

		It("calls to delete the instance", func() {
			instanceName := fakeInstanceDeleter.DeleteArgsForCall(0)
			Expect(instanceName).To(Equal("my-instance"))
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that the instance is being deleted", func() {
			Expect(fakeUI.DisplayTextCallCount()).To(Equal(1))

			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Instance '{{.InstanceName}}' is being deleted."))
			Expect(data[0]).To(HaveKeyWithValue("InstanceName", "my-instance"))
		})

		When("delete instance fails", func() {
			BeforeEach(func() {
				fakeInstanceDeleter.DeleteReturns(errors.New("error-deleting-instance"))
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError("error-deleting-instance"))
			})
		})
	})
})
