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

var _ = Describe("Binding List Usecase", func() {

	var (
		fakeBindingFetcher  *usecasesfakes.FakeBindingFetcher
		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher

		bindingListUsecase BindingListUsecase

		bindings   []*Binding
		executeErr error
	)

	BeforeEach(func() {
		fakeBindingFetcher = &usecasesfakes.FakeBindingFetcher{}
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}

		bindingListUsecase = BindingListUsecase{
			BindingFetcher:  fakeBindingFetcher,
			InstanceFetcher: fakeInstanceFetcher,
		}
	})

	JustBeforeEach(func() {
		bindings, executeErr = bindingListUsecase.GetBindings()
	})

	When("there are no bindings", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingsReturns([]*osbapi.Binding{}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns an empty list of bindings", func() {
			Expect(bindings).To(HaveLen(0))
		})
	})

	When("there are one or more bindings", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingsReturns([]*osbapi.Binding{
				{
					ID:         "binding-1",
					Name:       "my-binding-1",
					InstanceID: "instance-id-1",
					Status:     "creating",
					CreatedAt:  "time-1",
				}, {
					ID:         "binding-2",
					Name:       "my-binding-2",
					InstanceID: "instance-id-2",
					Status:     "created",
					CreatedAt:  "time-2",
				},
			}, nil)

			fakeInstanceFetcher.GetInstanceByIDReturnsOnCall(0, &osbapi.Instance{ID: "instance-id-1", Name: "my-instance-1"}, nil)
			fakeInstanceFetcher.GetInstanceByIDReturnsOnCall(1, &osbapi.Instance{ID: "instance-id-2", Name: "my-instance-2"}, nil)
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})

		It("returns the bindings", func() {
			Expect(bindings).To(HaveLen(2))

			Expect(*bindings[0]).To(Equal(Binding{
				Name:         "my-binding-1",
				InstanceName: "my-instance-1",
				Status:       "creating",
				CreatedAt:    "time-1",
			}))

			Expect(*bindings[1]).To(Equal(Binding{
				Name:         "my-binding-2",
				InstanceName: "my-instance-2",
				CreatedAt:    "time-2",
				Status:       "created",
			}))
		})

		It("fetches the instance using the instance id", func() {
			instanceID := fakeInstanceFetcher.GetInstanceByIDArgsForCall(0)
			Expect(instanceID).To(Equal("instance-id-1"))
		})
	})

	When("fetching bindings errors", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingsReturns([]*osbapi.Binding{}, errors.New("error-getting-bindings"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-bindings"))
		})
	})

	When("fetching the instance errors", func() {
		BeforeEach(func() {
			fakeBindingFetcher.GetBindingsReturns([]*osbapi.Binding{{
				ID:         "binding-1",
				Name:       "my-binding-1",
				ServiceID:  "service-id-1",
				PlanID:     "plan-id-1",
				BrokerName: "my-broker-1",
				Status:     "creating",
				CreatedAt:  "time-1",
			}}, nil)

			fakeInstanceFetcher.GetInstanceByIDReturns(&osbapi.Instance{}, errors.New("error-getting-instance"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-instance"))
		})
	})
})
