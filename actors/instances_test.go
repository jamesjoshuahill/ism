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

var _ = Describe("Instance Actor", func() {

	var (
		fakeInstanceRepository *actorsfakes.FakeInstanceRepository
		instancesActor         *InstancesActor
	)

	BeforeEach(func() {
		fakeInstanceRepository = &actorsfakes.FakeInstanceRepository{}

		instancesActor = &InstancesActor{
			Repository: fakeInstanceRepository,
		}
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			err = instancesActor.Create(&osbapi.Instance{
				Name: "instance-1",
			})
		})

		It("create the instance", func() {
			Expect(fakeInstanceRepository.CreateArgsForCall(0)).To(Equal(&osbapi.Instance{
				Name: "instance-1",
			}))
		})

		When("creating the instance fails", func() {
			BeforeEach(func() {
				fakeInstanceRepository.CreateReturns(errors.New("error-creating-instance"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-creating-instance"))
			})
		})
	})

	Describe("FindAll", func() {
		var (
			err       error
			instances []*osbapi.Instance
		)

		BeforeEach(func() {
			fakeInstanceRepository.FindAllReturns([]*osbapi.Instance{
				{
					Name: "my-instance-2",
				},
				{
					Name: "my-instance",
				},
			}, nil)
		})

		JustBeforeEach(func() {
			instances, err = instancesActor.GetInstances()
		})

		It("finds all the instances", func() {
			Expect(instances).To(HaveLen(2))
			Expect(instances).To(ConsistOf(
				&osbapi.Instance{Name: "my-instance"},
				&osbapi.Instance{Name: "my-instance-2"},
			))
		})

		When("finding the instances fails", func() {
			BeforeEach(func() {
				fakeInstanceRepository.FindAllReturns(nil, errors.New("error-finding-instances"))
			})

			It("propagates the error", func() {
				Expect(err).To(MatchError("error-finding-instances"))
			})
		})
	})
})
