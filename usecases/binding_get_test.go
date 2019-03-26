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

var _ = Describe("Binding Get Usecase", func() {
	var (
		bindingGetUsecase BindingGetUsecase

		fakeBindingFetcher  *usecasesfakes.FakeBindingFetcher
		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher
		fakeServiceFetcher  *usecasesfakes.FakeServiceFetcher
		fakePlanFetcher     *usecasesfakes.FakePlanFetcher
		fakeBrokerFetcher   *usecasesfakes.FakeBrokerFetcher

		bindingDetails *BindingDetails
		executeErr     error
	)

	BeforeEach(func() {
		fakeBindingFetcher = &usecasesfakes.FakeBindingFetcher{}
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}
		fakeServiceFetcher = &usecasesfakes.FakeServiceFetcher{}
		fakePlanFetcher = &usecasesfakes.FakePlanFetcher{}
		fakeBrokerFetcher = &usecasesfakes.FakeBrokerFetcher{}

		bindingGetUsecase = BindingGetUsecase{
			BindingFetcher:  fakeBindingFetcher,
			InstanceFetcher: fakeInstanceFetcher,
			ServiceFetcher:  fakeServiceFetcher,
			PlanFetcher:     fakePlanFetcher,
			BrokerFetcher:   fakeBrokerFetcher,
		}
	})

	JustBeforeEach(func() {
		bindingDetails, executeErr = bindingGetUsecase.GetBindingDetailsByName("my-binding")
	})

	When("the binding exists", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingByNameReturns(&osbapi.Binding{
				Name:        "my-binding",
				InstanceID:  "instance-1",
				ServiceID:   "service-1",
				PlanID:      "plan-1",
				BrokerName:  "my-broker",
				Status:      "creating",
				CreatedAt:   "time-1",
				Credentials: map[string]interface{}{"username": "admin"},
			}, nil)
			fakeInstanceFetcher.GetInstanceByIDReturns(&osbapi.Instance{
				ID:   "instance-1",
				Name: "my-instance",
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
			Expect(fakeBindingFetcher.GetBindingByNameArgsForCall(0)).To(Equal("my-binding"))
			Expect(fakeInstanceFetcher.GetInstanceByIDArgsForCall(0)).To(Equal("instance-1"))
			Expect(fakeServiceFetcher.GetServiceByIDArgsForCall(0)).To(Equal("service-1"))
			Expect(fakePlanFetcher.GetPlanByIDArgsForCall(0)).To(Equal("plan-1"))
		})

		It("returns the binding details", func() {
			Expect(*bindingDetails).To(Equal(BindingDetails{
				Name:         "my-binding",
				InstanceName: "my-instance",
				ServiceName:  "my-service",
				PlanName:     "my-plan",
				BrokerName:   "my-broker",
				Status:       "creating",
				CreatedAt:    "time-1",
				Credentials:  map[string]interface{}{"username": "admin"},
			}))
		})
	})

	Describe("errors", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingByNameReturns(&osbapi.Binding{
				Name:       "my-binding",
				InstanceID: "instance-1",
				ServiceID:  "service-1",
				PlanID:     "plan-1",
				BrokerName: "my-broker",
				Status:     "creating",
				CreatedAt:  "time-1",
			}, nil)
			fakeInstanceFetcher.GetInstanceByIDReturns(&osbapi.Instance{
				ID:   "instance-1",
				Name: "my-instance",
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

		When("fetching the binding fails", func() {
			BeforeEach(func() {
				fakeBindingFetcher.GetBindingByNameReturns(nil, errors.New("error-fetching-binding"))
			})

			It("propagates the error", func() {
				Expect(executeErr).To(MatchError("error-fetching-binding"))
			})
		})

		When("fetching the instance fails", func() {
			BeforeEach(func() {
				fakeInstanceFetcher.GetInstanceByIDReturns(nil, errors.New("error-fetching-instance"))
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
