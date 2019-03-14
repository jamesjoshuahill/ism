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

var _ = Describe("Instance List Usecase", func() {

	var (
		fakeInstancesFetcher *usecasesfakes.FakeInstancesFetcher
		fakeServiceFetcher   *usecasesfakes.FakeServiceFetcher
		fakePlanFetcher      *usecasesfakes.FakePlanFetcher

		instanceListUsecase InstanceListUsecase

		instances  []*Instance
		executeErr error
	)

	BeforeEach(func() {
		fakeInstancesFetcher = &usecasesfakes.FakeInstancesFetcher{}
		fakeServiceFetcher = &usecasesfakes.FakeServiceFetcher{}
		fakePlanFetcher = &usecasesfakes.FakePlanFetcher{}

		instanceListUsecase = InstanceListUsecase{
			InstancesFetcher: fakeInstancesFetcher,
			ServiceFetcher:   fakeServiceFetcher,
			PlanFetcher:      fakePlanFetcher,
		}
	})

	JustBeforeEach(func() {
		instances, executeErr = instanceListUsecase.GetInstances()
	})

	When("there are no instances", func() {
		BeforeEach(func() {
			fakeInstancesFetcher.GetInstancesReturns([]*osbapi.Instance{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns an empty list of services", func() {
			Expect(instances).To(HaveLen(0))
		})
	})

	When("there are one or more instances", func() {
		BeforeEach(func() {
			fakeInstancesFetcher.GetInstancesReturns([]*osbapi.Instance{
				{
					Name: "my-instance-1", ServiceID: "service-id-1", PlanID: "plan-id-1", BrokerName: "my-broker-1", CreatedAt: "time-1",
				}, {
					Name: "my-instance-2", ServiceID: "service-id-2", PlanID: "plan-id-2", BrokerName: "my-broker-2", CreatedAt: "time-2",
				},
			}, nil)

			fakeServiceFetcher.GetServiceReturnsOnCall(0, &osbapi.Service{Name: "my-service-1"}, nil)
			fakeServiceFetcher.GetServiceReturnsOnCall(1, &osbapi.Service{Name: "my-service-2"}, nil)
			fakePlanFetcher.GetPlanReturnsOnCall(0, &osbapi.Plan{Name: "my-plan-1"}, nil)
			fakePlanFetcher.GetPlanReturnsOnCall(1, &osbapi.Plan{Name: "my-plan-2"}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns the instances", func() {
			Expect(instances).To(HaveLen(2))

			Expect(*instances[0]).To(Equal(Instance{
				Name:        "my-instance-1",
				ServiceName: "my-service-1",
				PlanName:    "my-plan-1",
				BrokerName:  "my-broker-1",
				CreatedAt:   "time-1",
			}))

			Expect(*instances[1]).To(Equal(Instance{
				Name:        "my-instance-2",
				ServiceName: "my-service-2",
				PlanName:    "my-plan-2",
				BrokerName:  "my-broker-2",
				CreatedAt:   "time-2",
			}))
		})

		It("fetches the service using the service id", func() {
			serviceID := fakeServiceFetcher.GetServiceArgsForCall(0)
			Expect(serviceID).To(Equal("service-id-1"))
		})

		It("fetches the plan using the plan id", func() {
			planID := fakePlanFetcher.GetPlanArgsForCall(0)
			Expect(planID).To(Equal("plan-id-1"))
		})
	})

	When("fetching instances errors", func() {
		BeforeEach(func() {
			fakeInstancesFetcher.GetInstancesReturns([]*osbapi.Instance{}, errors.New("error-getting-instances"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-instances"))
		})
	})

	When("fetching the service errors", func() {
		BeforeEach(func() {
			fakeInstancesFetcher.GetInstancesReturns([]*osbapi.Instance{{
				Name:       "my-instance-1",
				ServiceID:  "service-id-1",
				PlanID:     "plan-id-1",
				BrokerName: "my-broker-1",
			}}, nil)

			fakeServiceFetcher.GetServiceReturns(&osbapi.Service{}, errors.New("error-getting-service"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-service"))
		})
	})

	When("fetching the plan errors", func() {
		BeforeEach(func() {
			fakeInstancesFetcher.GetInstancesReturns([]*osbapi.Instance{{
				Name:       "my-instance-1",
				ServiceID:  "service-id-1",
				PlanID:     "plan-id-1",
				BrokerName: "my-broker-1",
			}}, nil)

			fakeServiceFetcher.GetServiceReturns(&osbapi.Service{Name: "my-service-1"}, nil)
			fakePlanFetcher.GetPlanReturns(&osbapi.Plan{}, errors.New("error-getting-plan"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-plan"))
		})
	})
})
