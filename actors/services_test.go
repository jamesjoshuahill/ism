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

package actors_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/pivotal-cf/ism/actors"
	"github.com/pivotal-cf/ism/actors/actorsfakes"
	"github.com/pivotal-cf/ism/osbapi"
)

var _ = Describe("Services Actor", func() {

	var (
		fakeServiceRepository *actorsfakes.FakeServiceRepository

		servicesActor *ServicesActor
	)

	BeforeEach(func() {
		fakeServiceRepository = &actorsfakes.FakeServiceRepository{}

		servicesActor = &ServicesActor{
			Repository: fakeServiceRepository,
		}
	})

	Describe("GetService", func() {
		var (
			service *osbapi.Service
			err     error
		)

		BeforeEach(func() {
			fakeServiceRepository.FindReturns(&osbapi.Service{
				ID:   "service-id-1",
				Name: "service-1",
			}, nil)
		})

		JustBeforeEach(func() {
			service, err = servicesActor.GetService("service-id-1")
		})

		It("finds the service by service id", func() {
			Expect(fakeServiceRepository.FindArgsForCall(0)).To(Equal("service-id-1"))

			Expect(service).To(Equal(&osbapi.Service{
				Name: "service-1",
				ID:   "service-id-1",
			}))
		})

		When("finding the service returns an error", func() {
			BeforeEach(func() {
				fakeServiceRepository.FindReturns(nil, errors.New("error-finding-service"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-service"))
			})
		})
	})

	Describe("GetServices", func() {
		var (
			services []*osbapi.Service
			err      error
		)

		BeforeEach(func() {
			fakeServiceRepository.FindByBrokerReturns([]*osbapi.Service{
				{Name: "service-1"},
				{Name: "service-2"},
			}, nil)
		})

		JustBeforeEach(func() {
			services, err = servicesActor.GetServices("broker-1")
		})

		It("finds services by broker name", func() {
			Expect(fakeServiceRepository.FindByBrokerArgsForCall(0)).To(Equal("broker-1"))

			Expect(services).To(Equal([]*osbapi.Service{
				{Name: "service-1"},
				{Name: "service-2"},
			}))
		})

		When("finding services returns an error", func() {
			BeforeEach(func() {
				fakeServiceRepository.FindByBrokerReturns([]*osbapi.Service{}, errors.New("error-finding-services"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-services"))
			})
		})
	})
})
