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
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Broker List Command", func() {

	var (
		fakeBrokerFetcher *commandsfakes.FakeBrokerFetcher
		fakeUI            *commandsfakes.FakeUI

		listCommand BrokerListCommand

		executeErr error
	)

	BeforeEach(func() {
		fakeBrokerFetcher = &commandsfakes.FakeBrokerFetcher{}
		fakeUI = &commandsfakes.FakeUI{}

		listCommand = BrokerListCommand{
			BrokerFetcher: fakeBrokerFetcher,
			UI:            fakeUI,
		}
	})

	JustBeforeEach(func() {
		executeErr = listCommand.Execute(nil)
	})

	When("there are no brokers", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that no brokers were found", func() {
			Expect(fakeUI.DisplayTextCallCount()).NotTo(BeZero())
			text, _ := fakeUI.DisplayTextArgsForCall(0)

			Expect(text).To(Equal("No brokers found."))
		})
	})

	When("there is 1 or more brokers", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{{
				Name:      "broker-1",
				URL:       "https://broker-1-url.com",
				CreatedAt: "2019-02-28T12:08:31Z",
			}, {
				Name:      "broker-2",
				URL:       "https://broker-2-url.com",
				CreatedAt: "2018-02-27T12:09:30Z",
			}}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays the brokers in a table ", func() {
			Expect(fakeUI.DisplayTableCallCount()).NotTo(BeZero())
			data := fakeUI.DisplayTableArgsForCall(0)

			Expect(data[0]).To(Equal([]string{"NAME", "URL", "CREATED AT"}))
			Expect(data[1]).To(Equal([]string{"broker-1", "https://broker-1-url.com", "2019-02-28T12:08:31Z"}))
			Expect(data[2]).To(Equal([]string{"broker-2", "https://broker-2-url.com", "2018-02-27T12:09:30Z"}))
		})
	})

	When("getting brokers returns an error", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, errors.New("error-getting-brokers"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-brokers"))
		})
	})
})
