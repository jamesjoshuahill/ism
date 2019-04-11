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

	"github.com/pivotal-cf/ism/repositories"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/commands"
	"github.com/pivotal-cf/ism/commands/commandsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Broker command", func() {
	Describe("Register sub command", func() {
		var (
			fakeUI              *commandsfakes.FakeUI
			fakeBrokerRegistrar *commandsfakes.FakeBrokerRegistrar

			registerCommand BrokerRegisterCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeUI = &commandsfakes.FakeUI{}
			fakeBrokerRegistrar = &commandsfakes.FakeBrokerRegistrar{}

			registerCommand = BrokerRegisterCommand{
				UI:              fakeUI,
				BrokerRegistrar: fakeBrokerRegistrar,
			}
		})

		JustBeforeEach(func() {
			executeErr = registerCommand.Execute(nil)
		})

		When("given all required args", func() {
			BeforeEach(func() {
				registerCommand.Name = "broker-1"
				registerCommand.URL = "test-url"
				registerCommand.Username = "test-username"
				registerCommand.Password = "test-password"
			})

			It("doesn't error", func() {
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("displays that the broker was registered", func() {
				text, data := fakeUI.DisplayTextArgsForCall(0)
				Expect(text).To(Equal("Broker '{{.BrokerName}}' registered."))
				Expect(data[0]).To(HaveKeyWithValue("BrokerName", "broker-1"))
			})

			It("calls to register the broker", func() {
				broker := fakeBrokerRegistrar.RegisterArgsForCall(0)

				Expect(broker).To(Equal(&osbapi.Broker{
					Name:     "broker-1",
					URL:      "test-url",
					Username: "test-username",
					Password: "test-password",
				}))
			})

			When("a broker with the same name already exists", func() {
				BeforeEach(func() {
					fakeBrokerRegistrar.RegisterReturns(repositories.BrokerAlreadyExistsError)
				})

				It("returns an informative error message", func() {
					Expect(executeErr).To(MatchError("ERROR: A service broker named 'broker-1' already exists."))
				})
			})

			When("registering the broker errors", func() {
				BeforeEach(func() {
					fakeBrokerRegistrar.RegisterReturns(errors.New("error-registering-broker"))
				})

				It("propagates the error", func() {
					Expect(executeErr).To(MatchError("error-registering-broker"))
				})
			})
		})
	})

	Describe("List sub command", func() {
		var (
			fakeBrokersFetcher *commandsfakes.FakeBrokersFetcher
			fakeUI             *commandsfakes.FakeUI

			listCommand BrokerListCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeBrokersFetcher = &commandsfakes.FakeBrokersFetcher{}
			fakeUI = &commandsfakes.FakeUI{}

			listCommand = BrokerListCommand{
				BrokersFetcher: fakeBrokersFetcher,
				UI:             fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = listCommand.Execute(nil)
		})

		When("there are no brokers", func() {
			BeforeEach(func() {
				fakeBrokersFetcher.GetBrokersReturns([]*osbapi.Broker{}, nil)
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
				fakeBrokersFetcher.GetBrokersReturns([]*osbapi.Broker{{
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
				fakeBrokersFetcher.GetBrokersReturns([]*osbapi.Broker{}, errors.New("error-getting-brokers"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-getting-brokers"))
			})
		})
	})

	Describe("Delete sub command", func() {
		var (
			fakeBrokerDeleteUsecase *commandsfakes.FakeBrokerDeleteUsecase
			fakeUI                  *commandsfakes.FakeUI

			deleteCommand BrokerDeleteCommand

			executeErr error
		)

		BeforeEach(func() {
			fakeBrokerDeleteUsecase = &commandsfakes.FakeBrokerDeleteUsecase{}
			fakeUI = &commandsfakes.FakeUI{}

			deleteCommand = BrokerDeleteCommand{
				Name:                "my-broker",
				BrokerDeleteUsecase: fakeBrokerDeleteUsecase,
				UI:                  fakeUI,
			}
		})

		JustBeforeEach(func() {
			executeErr = deleteCommand.Execute(nil)
		})

		It("calls to delete the broker", func() {
			brokerName := fakeBrokerDeleteUsecase.DeleteArgsForCall(0)
			Expect(brokerName).To(Equal("my-broker"))
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("displays that the broker is being deleted", func() {
			Expect(fakeUI.DisplayTextCallCount()).To(Equal(1))

			text, data := fakeUI.DisplayTextArgsForCall(0)
			Expect(text).To(Equal("Broker '{{.BrokerName}}' is being deleted."))
			Expect(data[0]).To(HaveKeyWithValue("BrokerName", "my-broker"))
		})

		When("delete broker fails", func() {
			BeforeEach(func() {
				fakeBrokerDeleteUsecase.DeleteReturns(errors.New("error-deleting-broker"))
			})

			It("returns an error", func() {
				Expect(executeErr).To(MatchError("error-deleting-broker"))
			})
		})
	})
})
