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

var _ = Describe("Instance Delete Usecase", func() {
	var (
		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher
		fakeBindingFetcher  *usecasesfakes.FakeBindingFetcher
		fakeInstanceDeleter *usecasesfakes.FakeInstanceDeleter

		instanceDeleteUsecase InstanceDeleteUsecase

		instance *osbapi.Instance
		binding  *osbapi.Binding

		executeErr error
	)

	BeforeEach(func() {
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}
		fakeBindingFetcher = &usecasesfakes.FakeBindingFetcher{}
		fakeInstanceDeleter = &usecasesfakes.FakeInstanceDeleter{}

		instanceDeleteUsecase = InstanceDeleteUsecase{
			InstanceFetcher: fakeInstanceFetcher,
			BindingFetcher:  fakeBindingFetcher,
			InstanceDeleter: fakeInstanceDeleter,
		}

		instance = &osbapi.Instance{
			ID:   "instance-1",
			Name: "my-instance",
		}

		binding = &osbapi.Binding{
			ID:         "binding-1",
			Name:       "my-binding",
			InstanceID: "instance-1",
		}
	})

	JustBeforeEach(func() {
		executeErr = instanceDeleteUsecase.Delete("my-instance")
	})

	When("the instance does not have any bindings", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingFetcher.GetBindingsForInstanceReturns([]*osbapi.Binding{}, nil)
		})

		It("calls fetchers correctly", func() {
			Expect(fakeInstanceFetcher.GetInstanceByNameArgsForCall(0)).To(Equal("my-instance"))
			Expect(fakeBindingFetcher.GetBindingsForInstanceArgsForCall(0)).To(Equal("instance-1"))
		})

		It("deletes the instance", func() {
			Expect(fakeInstanceDeleter.DeleteCallCount()).To(Equal(1))
			Expect(fakeInstanceDeleter.DeleteArgsForCall(0)).To(Equal("my-instance"))
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})
	})

	When("the instance has one binding", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingFetcher.GetBindingsForInstanceReturns([]*osbapi.Binding{binding}, nil)
		})

		It("doesn't delete the instance", func() {
			Expect(fakeInstanceDeleter.DeleteCallCount()).To(Equal(0))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Instance 'my-instance' cannot be deleted as it has 1 binding"))
		})
	})

	When("the instance has more than one binding", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingFetcher.GetBindingsForInstanceReturns([]*osbapi.Binding{binding, &osbapi.Binding{}}, nil)
		})

		It("doesn't delete the instance", func() {
			Expect(fakeInstanceDeleter.DeleteCallCount()).To(Equal(0))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Instance 'my-instance' cannot be deleted as it has 2 bindings"))
		})
	})

	When("fetching the instance errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(&osbapi.Instance{}, errors.New("error-getting-instance"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-instance"))
		})
	})

	When("fetching the bindings for the instance errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingFetcher.GetBindingsForInstanceReturns([]*osbapi.Binding{}, errors.New("error-getting-bindings"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-bindings"))
		})
	})

	When("deleting the instance errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstanceByNameReturns(instance, nil)
			fakeBindingFetcher.GetBindingsForInstanceReturns([]*osbapi.Binding{}, nil)
			fakeInstanceDeleter.DeleteReturns(errors.New("error-deleting-instance"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-deleting-instance"))
		})
	})
})
