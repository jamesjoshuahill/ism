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

var _ = Describe("Instance List Command", func() {

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

	When("there is 1 or more instance", func() {
		BeforeEach(func() {
			fakeInstanceListUsecase.GetInstancesReturns([]*usecases.Instance{{
				Name:        "my-instance-1",
				ServiceName: "my-service-1",
				PlanName:    "my-plan-1",
				BrokerName:  "my-broker-1",
				CreatedAt:   "2019-02-28T12:08:31Z",
			}, {
				Name:        "my-instance-2",
				ServiceName: "my-service-1",
				PlanName:    "my-plan-2",
				BrokerName:  "my-broker-1",
				CreatedAt:   "2019-02-28T12:08:31Z",
			}}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays the instances in a table ", func() {
			Expect(fakeUI.DisplayTableCallCount()).NotTo(BeZero())
			data := fakeUI.DisplayTableArgsForCall(0)

			Expect(data[0]).To(Equal([]string{"NAME", "SERVICE", "PLAN", "BROKER", "CREATED AT"}))
			Expect(data[1]).To(Equal([]string{"my-instance-1", "my-service-1", "my-plan-1", "my-broker-1", "2019-02-28T12:08:31Z"}))
			Expect(data[2]).To(Equal([]string{"my-instance-2", "my-service-1", "my-plan-2", "my-broker-1", "2019-02-28T12:08:31Z"}))
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
