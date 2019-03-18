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

var _ = Describe("Binding Create Usecase", func() {
	var (
		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher
		fakeBindingCreator  *usecasesfakes.FakeBindingCreator

		bindingCreateUsecase BindingCreateUsecase

		instance *osbapi.Instance

		executeErr error
	)

	BeforeEach(func() {
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}
		fakeBindingCreator = &usecasesfakes.FakeBindingCreator{}

		instance = &osbapi.Instance{
			ID:         "instance-1",
			Name:       "my-instance-1",
			PlanID:     "plan-1",
			ServiceID:  "service-1",
			BrokerName: "my-broker-1",
		}

		bindingCreateUsecase = BindingCreateUsecase{
			BindingCreator:  fakeBindingCreator,
			InstanceFetcher: fakeInstanceFetcher,
		}
	})

	JustBeforeEach(func() {
		executeErr = bindingCreateUsecase.Create("my-binding-1", "my-instance-1")
	})

	When("the instance exists", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
		})

		It("creates a binding", func() {
			Expect(fakeInstanceFetcher.GetInstanceByNameCallCount()).To(Equal(1))

			passedInstanceName := fakeInstanceFetcher.GetInstanceByNameArgsForCall(0)
			Expect(passedInstanceName).To(Equal("my-instance-1"))

			Expect(fakeBindingCreator.CreateCallCount()).To(Equal(1))
			passedBinding := fakeBindingCreator.CreateArgsForCall(0)

			Expect(*passedBinding).To(Equal(
				osbapi.Binding{
					Name:       "my-binding-1",
					InstanceID: "instance-1",
					PlanID:     "plan-1",
					ServiceID:  "service-1",
					BrokerName: "my-broker-1",
				},
			))

			Expect(executeErr).NotTo(HaveOccurred())
		})
	})

	When("the instance does not exist", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(nil, nil)
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Instance 'my-instance-1' does not exist"))
		})
	})

	When("fetching the instance errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(nil, errors.New("error-fetching-instance"))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("error-fetching-instance"))
		})
	})

	When("creating the instance errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingCreator.CreateReturns(errors.New("error-creating-binding"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-creating-binding"))
		})
	})
})
