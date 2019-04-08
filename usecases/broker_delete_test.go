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

var _ = Describe("Broker Delete Usecase", func() {
	var (
		fakeInstanceFetcher *usecasesfakes.FakeInstanceFetcher
		fakeBrokerDeleter   *usecasesfakes.FakeBrokerDeleter

		brokerDeleteUsecase BrokerDeleteUsecase

		instance *osbapi.Instance

		executeErr error
	)

	BeforeEach(func() {
		fakeInstanceFetcher = &usecasesfakes.FakeInstanceFetcher{}
		fakeBrokerDeleter = &usecasesfakes.FakeBrokerDeleter{}

		brokerDeleteUsecase = BrokerDeleteUsecase{
			InstanceFetcher: fakeInstanceFetcher,
			BrokerDeleter:   fakeBrokerDeleter,
		}

		instance = &osbapi.Instance{
			ID:         "instance-1",
			Name:       "my-instance",
			BrokerName: "my-broker",
		}
	})

	JustBeforeEach(func() {
		executeErr = brokerDeleteUsecase.Delete("my-broker")
	})

	When("the broker does not have any instances", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstancesForBrokerReturns([]*osbapi.Instance{}, nil)
		})

		It("calls fetchers correctly", func() {
			Expect(fakeInstanceFetcher.GetInstancesForBrokerArgsForCall(0)).To(Equal("my-broker"))
		})

		It("deletes the instance", func() {
			Expect(fakeBrokerDeleter.DeleteCallCount()).To(Equal(1))
			Expect(fakeBrokerDeleter.DeleteArgsForCall(0)).To(Equal("my-broker"))
		})

		It("doesn't error", func() {
			Expect(executeErr).NotTo(HaveOccurred())
		})
	})

	When("the broker has one instance", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstancesForBrokerReturns([]*osbapi.Instance{instance}, nil)
		})

		It("doesn't delete the broker", func() {
			Expect(fakeBrokerDeleter.DeleteCallCount()).To(Equal(0))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Broker 'my-broker' cannot be deleted as it has 1 instance"))
		})
	})

	When("the broker has more than one instance", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstancesForBrokerReturns([]*osbapi.Instance{instance, &osbapi.Instance{}}, nil)
		})

		It("doesn't delete the instance", func() {
			Expect(fakeBrokerDeleter.DeleteCallCount()).To(Equal(0))
		})

		It("returns an error", func() {
			Expect(executeErr).To(MatchError("Broker 'my-broker' cannot be deleted as it has 2 instances"))
		})
	})

	When("fetching the instances errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstancesForBrokerReturns([]*osbapi.Instance{}, errors.New("error-getting-instances"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-getting-instances"))
		})
	})

	When("deleting the broker errors", func() {
		BeforeEach(func() {
			fakeInstanceFetcher.GetInstancesForBrokerReturns([]*osbapi.Instance{}, nil)
			fakeBrokerDeleter.DeleteReturns(errors.New("error-deleting-instance"))
		})

		It("propagates the error", func() {
			Expect(executeErr).To(MatchError("error-deleting-instance"))
		})
	})
})
