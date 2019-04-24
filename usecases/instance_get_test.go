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

var _ = Describe("Instance Get Usecase", func() {
	var (
		instanceGetUsecase InstanceGetUsecase

		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher
		fakeServiceFetcher  *usecasesfakes.FakeServiceFetcher
		fakePlanFetcher     *usecasesfakes.FakePlanFetcher
		fakeBrokerFetcher   *usecasesfakes.FakeBrokerFetcher

		instanceDetails *InstanceDetails
		executeErr      error
	)

	BeforeEach(func() {
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}
		fakeServiceFetcher = &usecasesfakes.FakeServiceFetcher{}
		fakePlanFetcher = &usecasesfakes.FakePlanFetcher{}
		fakeBrokerFetcher = &usecasesfakes.FakeBrokerFetcher{}

		instanceGetUsecase = InstanceGetUsecase{
			InstanceFetcher: fakeInstanceFetcher,
			ServiceFetcher:  fakeServiceFetcher,
			PlanFetcher:     fakePlanFetcher,
			BrokerFetcher:   fakeBrokerFetcher,
		}
	})

	JustBeforeEach(func() {
		instanceDetails, executeErr = instanceGetUsecase.GetInstanceDetailsByName("my-instance")
	})

	When("the instance exists", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(&osbapi.Instance{
				ID:         "instance-1",
				Name:       "my-instance",
				CreatedAt:  "time-1",
				Status:     "created",
				ServiceID:  "service-1",
				PlanID:     "plan-1",
				BrokerName: "my-broker",
			}, nil)
			fakeServiceFetcher.GetServiceByIDReturns(&osbapi.Service{
				ID:   "service-1",
				Name: "my-service",
			}, nil)
			fakePlanFetcher.GetPlanByIDReturns(&osbapi.Plan{
				ID:   "plan-1",
				Name: "my-plan",
			}, nil)
			fakeBrokerFetcher.GetBrokerByNameReturns(&osbapi.Broker{
				Name: "my-broker",
			}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("calls fetchers correctly", func() {
			Expect(fakeInstanceFetcher.GetInstanceByNameArgsForCall(0)).To(Equal("my-instance"))
			Expect(fakeServiceFetcher.GetServiceByIDArgsForCall(0)).To(Equal("service-1"))
			Expect(fakePlanFetcher.GetPlanByIDArgsForCall(0)).To(Equal("plan-1"))
			Expect(fakePlanFetcher.GetPlanByIDArgsForCall(0)).To(Equal("plan-1"))
			Expect(fakeBrokerFetcher.GetBrokerByNameArgsForCall(0)).To(Equal("my-broker"))
		})

		It("returns the instance details", func() {
			Expect(*instanceDetails).To(Equal(InstanceDetails{
				Name:        "my-instance",
				CreatedAt:   "time-1",
				Status:      "created",
				ServiceName: "my-service",
				PlanName:    "my-plan",
				BrokerName:  "my-broker",
			}))
		})
	})

	Describe("errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(&osbapi.Instance{
				ID:         "instance-1",
				Name:       "my-instance",
				CreatedAt:  "time-1",
				Status:     "created",
				ServiceID:  "service-1",
				PlanID:     "plan-1",
				BrokerName: "my-broker",
			}, nil)
			fakeServiceFetcher.GetServiceByIDReturns(&osbapi.Service{
				ID:   "service-1",
				Name: "my-service",
			}, nil)
			fakePlanFetcher.GetPlanByIDReturns(&osbapi.Plan{
				ID:   "plan-1",
				Name: "my-plan",
			}, nil)
			fakeBrokerFetcher.GetBrokerByNameReturns(&osbapi.Broker{
				Name: "my-broker",
			}, nil)
		})

		When("fetching the instance fails", func() {
			BeforeEach(func() {
				fakeInstanceFetcher.GetInstanceByNameReturns(nil, errors.New("error-fetching-instance"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-fetching-instance"))
			})
		})

		When("fetching the service fails", func() {
			BeforeEach(func() {
				fakeServiceFetcher.GetServiceByIDReturns(nil, errors.New("error-fetching-service"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-fetching-service"))
			})
		})

		When("fetching the plan fails", func() {
			BeforeEach(func() {
				fakePlanFetcher.GetPlanByIDReturns(nil, errors.New("error-fetching-plan"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-fetching-plan"))
			})
		})

		When("fetching the broker fails", func() {
			BeforeEach(func() {
				fakeBrokerFetcher.GetBrokerByNameReturns(nil, errors.New("error-fetching-broker"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-fetching-broker"))
			})
		})
	})
})
