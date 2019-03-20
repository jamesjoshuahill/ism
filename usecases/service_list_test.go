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

package usecases_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/ism/osbapi"
	. "github.com/pivotal-cf/ism/usecases"
	"github.com/pivotal-cf/ism/usecases/usecasesfakes"
)

var _ = Describe("Service List Usecase", func() {

	var (
		fakeBrokerFetcher  *usecasesfakes.FakeBrokerFetcher
		fakeServiceFetcher *usecasesfakes.FakeServiceFetcher
		fakePlanFetcher    *usecasesfakes.FakePlanFetcher

		serviceListUsecase ServiceListUsecase

		services   []*Service
		executeErr error
	)

	BeforeEach(func() {
		fakeBrokerFetcher = &usecasesfakes.FakeBrokerFetcher{}
		fakeServiceFetcher = &usecasesfakes.FakeServiceFetcher{}
		fakePlanFetcher = &usecasesfakes.FakePlanFetcher{}

		serviceListUsecase = ServiceListUsecase{
			BrokerFetcher:  fakeBrokerFetcher,
			ServiceFetcher: fakeServiceFetcher,
			PlanFetcher:    fakePlanFetcher,
		}
	})

	JustBeforeEach(func() {
		services, executeErr = serviceListUsecase.GetServices()
	})

	It("fetches all brokers", func() {
		Expect(fakeBrokerFetcher.GetBrokersCallCount()).NotTo(BeZero())
	})

	When("fetching brokers errors", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, errors.New("error-getting-brokers"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-brokers"))
		})
	})

	When("there are no brokers", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns an empty list of services", func() {
			Expect(services).To(HaveLen(0))
		})
	})

	When("there are one or more brokers", func() {
		BeforeEach(func() {
			fakeBrokerFetcher.GetBrokersReturns([]*osbapi.Broker{
				{Name: "broker1"},
				{Name: "broker2"}}, nil)
		})

		It("fetches services for each broker", func() {
			Expect(fakeServiceFetcher.GetServicesCallCount()).To(Equal(2))
			Expect(fakeServiceFetcher.GetServicesArgsForCall(0)).To(Equal("broker1"))
			Expect(fakeServiceFetcher.GetServicesArgsForCall(1)).To(Equal("broker2"))
		})

		When("fetching services errors", func() {
			BeforeEach(func() {
				fakeServiceFetcher.GetServicesReturns([]*osbapi.Service{}, errors.New("error-getting-services"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-getting-services"))
			})
		})

		When("all the brokers have services", func() {
			BeforeEach(func() {
				fakeServiceFetcher.GetServicesReturnsOnCall(0, []*osbapi.Service{
					{ID: "service1-id", Name: "service1", Description: "service1 description"},
					{ID: "service2-id", Name: "service2", Description: "service2 description"}}, nil)

				fakeServiceFetcher.GetServicesReturnsOnCall(1, []*osbapi.Service{
					{ID: "service3-id", Name: "service3", Description: "service3 description"},
					{ID: "service4-id", Name: "service4", Description: "service4 description"}}, nil)
			})

			It("fetches plans for each service", func() {
				Expect(fakePlanFetcher.GetPlansCallCount()).To(Equal(4))
				Expect(fakePlanFetcher.GetPlansArgsForCall(0)).To(Equal("service1-id"))
				Expect(fakePlanFetcher.GetPlansArgsForCall(1)).To(Equal("service2-id"))
				Expect(fakePlanFetcher.GetPlansArgsForCall(2)).To(Equal("service3-id"))
				Expect(fakePlanFetcher.GetPlansArgsForCall(3)).To(Equal("service4-id"))
			})

			When("fetching plans errors", func() {
				BeforeEach(func() {
					fakePlanFetcher.GetPlansReturns([]*osbapi.Plan{}, errors.New("error-getting-plans"))
				})

				It("propagates the error", func() {
					Expect(executeErr).To(MatchError("error-getting-plans"))
				})
			})

			When("all the services have plans", func() {
				BeforeEach(func() {
					fakePlanFetcher.GetPlansReturnsOnCall(0, []*osbapi.Plan{{Name: "plan1"}, {Name: "extra-plan"}}, nil)
					fakePlanFetcher.GetPlansReturnsOnCall(1, []*osbapi.Plan{{Name: "plan2"}}, nil)
					fakePlanFetcher.GetPlansReturnsOnCall(2, []*osbapi.Plan{{Name: "plan3"}}, nil)
					fakePlanFetcher.GetPlansReturnsOnCall(3, []*osbapi.Plan{{Name: "plan4"}}, nil)
				})

				It("doesn't error", func() {
					Expect(executeErr).NotTo(HaveOccurred())
				})

				It("returns a list of services", func() {
					Expect(services[0]).To(Equal(&Service{Name: "service1", Description: "service1 description", PlanNames: []string{"plan1", "extra-plan"}, BrokerName: "broker1"}))
					Expect(services[1]).To(Equal(&Service{Name: "service2", Description: "service2 description", PlanNames: []string{"plan2"}, BrokerName: "broker1"}))
					Expect(services[2]).To(Equal(&Service{Name: "service3", Description: "service3 description", PlanNames: []string{"plan3"}, BrokerName: "broker2"}))
					Expect(services[3]).To(Equal(&Service{Name: "service4", Description: "service4 description", PlanNames: []string{"plan4"}, BrokerName: "broker2"}))
				})
			})
		})
	})
})
